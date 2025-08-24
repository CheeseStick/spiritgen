package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	// Convert to rendered structure
	rendered := render.FromSpiritTablets(result.Success)

	// Determine output path
	dir := filepath.Dir(*inputPath)
	outputPath := filepath.Join(dir, *outputName)

	// Generate PDF
	err = render.RenderLabelsAsPDF(rendered, outputPath)
	if err != nil {
		log.Fatalf("❌ PDF 생성 실패: %v", err)
	}

	fmt.Printf("✅ PDF generated: %s\n", outputPath)
}
