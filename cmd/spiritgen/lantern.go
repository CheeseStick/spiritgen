package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"spiritgen/assets"
	"spiritgen/internal/lantern"
)

func runLantern(args []string) {
	fs := flag.NewFlagSet("lantern", flag.ExitOnError)
	inputFlag := fs.String("input", "", "XLSX 파일 경로")
	outputName := fs.String("output", "lantern_output.pdf", "PDF 출력 파일 이름 (기본: lantern_output.pdf)")
	titleFlag := fs.String("title", "", "행사명 (모든 tablet 상단에 동일하게 출력. 비우면 '뉴질랜드 남국선사')")
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	inputPath := *inputFlag
	if inputPath == "" && len(fs.Args()) > 0 {
		inputPath = fs.Args()[0]
	}

	if inputPath == "" {
		log.Fatal("❌ XLSX 파일이 필요합니다. --input <path>")
	}
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		log.Fatalf("❌ 파일이 존재하지 않습니다: %s", inputPath)
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("❌ 데이터를 읽을 수 없습니다: %v", err)
	}

	result, err := lantern.ParseXLSX(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("❌ XLSX 파일을 로드하는데 실패했습니다: %v", err)
	}

	// Log per-row errors (per sheet) before deciding whether to proceed.
	totalErrs := result.TotalErrorCount()
	for _, sheetName := range []string{lantern.SheetBig, lantern.SheetFamily, lantern.SheetSpirit, lantern.SheetBusiness} {
		for _, rowErr := range result.Errors[sheetName] {
			for _, e := range rowErr.Errors {
				log.Printf("⚠️ [%s] 행 %d: %s [%s]", sheetName, rowErr.RowIndex, e.Message, e.Code)
			}
		}
	}
	if totalErrs > 0 {
		log.Printf("⚠️ 일부 행에서 유효성 오류가 발생했습니다 (%d개). 무시하고 진행합니다.", totalErrs)
	}

	totalHouseholds := result.TotalHouseholdCount()
	if totalHouseholds == 0 {
		log.Fatalf("❌ 처리 할 세대 데이터가 없습니다.")
	}

	outputPath := *outputName
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(filepath.Dir(inputPath), outputPath)
	}

	if err := lantern.RenderPDF(result, outputPath, assets.LanternTabletDesignOne, *titleFlag); err != nil {
		log.Fatalf("❌ PDF 생성 실패: %v", err)
	}

	fmt.Printf("✅ PDF가 성공적으로 생성되었습니다: %s (큰등 %d · 가족등 %d · 영가등 %d · 사업등 %d)\n",
		outputPath, len(result.Big), len(result.Family), len(result.Spirit), len(result.Business))
}
