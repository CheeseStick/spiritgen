package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"spiritgen/internal/model"
	"spiritgen/internal/parser"
	"spiritgen/internal/render"
)

func main() {
	inputFlag := flag.String("input", "", "XLSX 파일 경로 (또는 파일을 실행 파일 위로 드래그 하세요)")
	outputName := flag.String("output", "output.pdf", "PDF 출력 파일 이름 (기본: output.pdf)")
	flag.Parse()

	inputPath := *inputFlag
	if inputPath == "" && len(flag.Args()) > 0 {
		inputPath = flag.Args()[0]
	}

	if inputPath == "" {
		log.Fatal("❌ XLSX 파일이 필요합니다. --input <path> 또는 파일을 실행 파일 위로 드레 해주세요.")
	}

	// Validate input file
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		log.Fatalf("❌ 파일이 존재하지 않습니다.: %s", inputPath)
	}

	// Read input XLSX
	data, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("❌ 데이터를 읽을 수 없습니다: %v", err)
	}

	// Parse XLSX to SpiritTablets
	result, err := parser.ParseFromXLSX(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("❌ XLSX 파일을 로드하는데 실패했습니다: %v", err)
	}

	if 0 < len(result.Errors) {
		log.Printf("⚠️ 일부 행에서 유효성 오류가 발생했습니다 (%d개). 무시하고 진행합니다.", len(result.Errors))
	}

	dir := filepath.Dir(inputPath)
	outputPath := filepath.Join(dir, *outputName)

	// Split spirit tablets
	spiritTablets := make([]model.SpiritTablet, 0, len(result.Success))

	for _, tablet := range result.Success {
		spiritTablets = append(spiritTablets, tablet.Split(3)...)
	}

	if len(spiritTablets) == 0 {
		log.Fatalf("❌ 처리 할 데이터가 없습니다.")
	}

	// Generate PDF
	err = render.FromSpiritTablets(result.Success, outputPath)
	if err != nil {
		log.Fatalf("❌ PDF 생성 실패: %v", err)
	}

	fmt.Printf("✅ PDF가 성공적으로 생성되었습니다: %s\n", outputPath)
}
