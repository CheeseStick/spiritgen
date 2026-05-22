package lantern

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spiritgen/assets"
	"spiritgen/internal/gui"
	"spiritgen/internal/validation"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

const (
	designOptionOne = "디자인 1"
	designOptionTwo = "디자인 2"
)

// uiState holds the runtime state of the lantern tab.
type uiState struct {
	xlsxPath       string
	parseResult    *ParseResult
	households     []Household
	householdCount int
	personCount    int
	selectedDesign []byte
	title          string
	logPane        *gui.LogPane // nil-safe (e.g. in tests)
}

type tab struct {
	state *uiState
}

// NewTab returns a freshly-initialized lantern tab implementing gui.Tab.
func NewTab() gui.Tab {
	return &tab{
		state: &uiState{selectedDesign: assets.LanternTabletDesignOne},
	}
}

func (t *tab) Title() string { return "인등/연등 생성" }

func (t *tab) Build(win fyne.Window) fyne.CanvasObject {
	state := t.state

	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("xlsx 파일을 선택해주세요...")
	pathEntry.Disable()

	statusLabel := widget.NewLabel("")
	statusLabel.Alignment = fyne.TextAlignCenter
	statusLabel.Hide()

	generateBtn := widget.NewButton("생성", nil)
	generateBtn.Importance = widget.HighImportance
	generateBtn.Disable()

	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()

	afterRead := func() {
		if state.householdCount > 0 {
			statusLabel.SetText(fmt.Sprintf("✓  %d세대 · %d명 로드됨", state.householdCount, state.personCount))
			statusLabel.Show()
		}
	}

	openFilePicker := func(onSelected func()) {
		fd := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if err != nil || r == nil {
				return
			}
			defer r.Close()
			state.xlsxPath = r.URI().Path()
			pathEntry.SetText(state.xlsxPath)
			if onSelected != nil {
				onSelected()
			}
		}, win)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		fd.Show()
	}

	readBtn := widget.NewButton("읽기", func() {
		if state.xlsxPath == "" {
			openFilePicker(func() {
				DoRead(state, generateBtn, win)
				afterRead()
			})
			return
		}
		DoRead(state, generateBtn, win)
		afterRead()
	})

	browseBtn := widget.NewButton("파일 선택", func() {
		openFilePicker(nil)
	})

	templateBtn := widget.NewButton("템플릿 저장", func() {
		fd := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
			if err != nil || w == nil {
				return
			}
			defer w.Close()
			if _, writeErr := w.Write(assets.LanternTabletTemplate); writeErr != nil {
				dialog.ShowError(writeErr, win)
				return
			}
			dialog.ShowInformation("완료", "템플릿 파일이 저장되었습니다.", win)
		}, win)
		fd.SetFileName("lantern_tablet_template.xlsx")
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		fd.Show()
	})
	templateBtn.Importance = widget.LowImportance

	designSelector := widget.NewRadioGroup(
		[]string{designOptionOne, designOptionTwo},
		func(selected string) {
			if selected == designOptionTwo {
				state.selectedDesign = assets.LanternTabletDesignTwo
			} else {
				state.selectedDesign = assets.LanternTabletDesignOne
			}
		},
	)
	designSelector.SetSelected(designOptionOne)
	designSelector.Horizontal = true

	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("예: 동지기도, 부처님오신날 (비워두면 '뉴질랜드 남국선사')")
	titleEntry.OnChanged = func(s string) { state.title = s }

	generateBtn.OnTapped = func() {
		generateBtn.Disable()
		readBtn.Disable()
		progressBar.Show()
		progressBar.Start()
		showSaveDialog(state, readBtn, generateBtn, progressBar, win)
	}

	state.logPane = gui.NewLogPane()

	fileRow := container.NewBorder(nil, nil, nil, browseBtn, pathEntry)
	form := widget.NewForm(
		widget.NewFormItem("파일", fileRow),
		widget.NewFormItem("행사명", titleEntry),
		widget.NewFormItem("디자인", designSelector),
		widget.NewFormItem("", container.NewHBox(templateBtn)),
	)

	actionRow := container.NewGridWithColumns(2, readBtn, generateBtn)

	content := container.NewVBox(
		form,
		widget.NewSeparator(),
		statusLabel,
		progressBar,
		actionRow,
		state.logPane.Widget(),
	)

	return container.NewPadded(content)
}

func showSaveDialog(state *uiState, readBtn, generateBtn *widget.Button, progressBar *widget.ProgressBarInfinite, win fyne.Window) {
	restoreButtons := func() {
		progressBar.Stop()
		progressBar.Hide()
		generateBtn.Enable()
		readBtn.Enable()
	}

	filenameEntry := widget.NewEntry()
	filenameEntry.SetText("lantern_output.pdf")

	folderLabel := widget.NewLabel("폴더를 선택해주세요")
	var chosenFolder string

	folderBtn := widget.NewButton("선택...", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			chosenFolder = uri.Path()
			folderLabel.SetText(chosenFolder)
		}, win).Show()
	})

	items := []*widget.FormItem{
		widget.NewFormItem("저장 폴더", container.NewBorder(nil, nil, nil, folderBtn, folderLabel)),
		widget.NewFormItem("파일 이름", filenameEntry),
	}

	dialog.ShowForm("저장 위치 선택", "저장", "취소", items, func(confirmed bool) {
		if !confirmed {
			restoreButtons()
			return
		}
		if chosenFolder == "" {
			dialog.ShowError(fmt.Errorf("저장 폴더를 선택해주세요"), win)
			restoreButtons()
			return
		}

		filename := strings.TrimSpace(filenameEntry.Text)
		if filename == "" {
			filename = "lantern_output.pdf"
		}
		if !strings.HasSuffix(strings.ToLower(filename), ".pdf") {
			filename += ".pdf"
		}

		outputPath := filepath.Join(chosenFolder, filename)
		design := state.selectedDesign
		title := state.title

		doRender := func() {
			go func() {
				renderErr := RenderPDF(state.households, outputPath, design, title)
				restoreButtons()
				if renderErr != nil {
					dialog.ShowError(renderErr, win)
				} else {
					dialog.ShowInformation("완료", "PDF가 생성되었습니다:\n"+outputPath, win)
				}
			}()
		}

		if _, err := os.Stat(outputPath); err == nil {
			dialog.ShowConfirm(
				"덮어쓰기 확인",
				fmt.Sprintf("'%s' 파일이 이미 존재합니다.\n덮어쓰시겠습니까?", filename),
				func(overwrite bool) {
					if !overwrite {
						restoreButtons()
						return
					}
					doRender()
				}, win)
			return
		}

		doRender()
	}, win)
}

// DoRead parses the xlsx at state.xlsxPath and updates state. Exported for
// testing — returns without action if xlsxPath is empty.
func DoRead(state *uiState, generateBtn *widget.Button, win fyne.Window) {
	if state.xlsxPath == "" {
		return
	}

	// Reset the log pane for a fresh read.
	state.logPane.Clear()
	state.logPane.Append(fmt.Sprintf("📂 %s", state.xlsxPath))

	f, err := os.Open(state.xlsxPath)
	if err != nil {
		state.logPane.Append(fmt.Sprintf("❌ 파일을 열 수 없습니다: %v", err))
		dialog.ShowError(fmt.Errorf("파일을 열 수 없습니다: %w", err), win)
		return
	}
	defer f.Close()

	result, err := ParseXLSX(f)
	if err != nil {
		state.logPane.Append(fmt.Sprintf("❌ XLSX 파싱 실패: %v", err))
		dialog.ShowError(fmt.Errorf("XLSX 파싱 실패: %w", err), win)
		return
	}

	// Always log per-row errors first — so the user can see WHY rows failed
	// even when nothing valid was parsed.
	for _, rowErr := range result.Errors {
		for _, e := range rowErr.Errors {
			state.logPane.Append(fmt.Sprintf("⚠️ 행 %d: %s [%s]", rowErr.RowIndex, e.Message, e.Code))
		}
	}

	if len(result.Success) == 0 {
		state.logPane.Append(fmt.Sprintf("❌ 유효한 데이터가 없습니다 (오류 %d행)", len(result.Errors)))
		dialog.ShowError(fmt.Errorf("유효한 데이터가 없습니다 (오류 %d행)", len(result.Errors)), win)
		return
	}

	personCount := 0
	for _, h := range result.Success {
		personCount += len(h.Persons)
	}

	state.parseResult = &result
	state.households = result.Success
	state.householdCount = len(result.Success)
	state.personCount = personCount

	state.logPane.Append(fmt.Sprintf("✓ %d세대 · %d명 로드 완료 (오류 %d행)",
		state.householdCount, personCount, len(result.Errors)))

	if len(result.Errors) > 0 {
		showRowErrorWarning(result.Errors, win)
	}

	generateBtn.Enable()
	dialog.ShowInformation("읽기 완료", fmt.Sprintf("%d세대 · %d명을 읽었습니다.", state.householdCount, personCount), win)
}

func showRowErrorWarning(errs []validation.RowError, win fyne.Window) {
	dialog.ShowInformation("경고",
		fmt.Sprintf("%d개 행에서 오류가 발생했습니다. 무시하고 진행합니다.", len(errs)),
		win)
}
