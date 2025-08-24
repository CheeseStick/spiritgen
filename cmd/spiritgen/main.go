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
	// Define flags
	inputPath := flag.String("input", "", "Path to XLSX input file")
	outputName := flag.String("output", "output.pdf", "Name of output PDF file (optional)")

	flag.Parse()

	if *inputPath == "" {
		log.Fatal("❌ input XLSX file path is required. Use --input <path>")
	}

	// Validate input file
	if _, err := os.Stat(*inputPath); os.IsNotExist(err) {
		log.Fatalf("❌ input file does not exist: %s", *inputPath)
	}

	// Read input XLSX
	data, err := os.ReadFile(*inputPath)
	if err != nil {
		log.Fatalf("❌ failed to read input file: %v", err)
	}

	// Parse XLSX to SpiritTablets
	result, err := parser.ParseFromXLSX(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("❌ failed to parse XLSX: %v", err)
	}

	if 0 < len(result.Errors) {
		log.Printf("⚠️ 일부 행에서 유효성 오류가 발생했습니다 (%d개). 무시하고 진행합니다.", len(result.Errors))
	}

	// Determine output path
	dir := filepath.Dir(*inputPath)
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

	fmt.Printf("✅ PDF generated: %s\n", outputPath)
}
