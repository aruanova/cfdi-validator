// cfdi-validator/main.go
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

type CFDIDetails2022 struct {
	CfdiId                              string          `json:"cfdiId"`
	Status                              string          `json:"status"`
	CfdiRelationType                    string          `json:"cfdiRelationType"`
	RelatedCfdi                         []uuid.UUID     `json:"relatedCfdi,omitempty"`
	Type                                string          `json:"type"`
	Uuid                                uuid.UUID       `json:"uuid"`
	Series                              string          `json:"series"`
	Reference                           string          `json:"reference"`
	EmitterRfc                          string          `json:"emitterRfc"`
	EmitterCompanyName                  string          `json:"emitterCompanyName"`
	EmitterPostalCode                   string          `json:"emitterPostalCode"`
	CfdiUsage                           string          `json:"cfdiUsage"`
	ReceiverRfc                         string          `json:"receiverRfc"`
	ReceiverCompanyName                 string          `json:"receiverCompanyName"`
	ReceptorPostalCode                  string          `json:"receptorPostalCode"`
	Origin                              string          `json:"origin"`
	Currency                            string          `json:"currency"`
	StampedDate                         string          `json:"stampedDate"`
	InvoiceDate                         string          `json:"invoiceDate"`
	Grouping                            string          `json:"grouping"`
	Iva                                 decimal.Decimal `json:"iva"`
	SubTotal                            decimal.Decimal `json:"subTotal"`
	Total                               decimal.Decimal `json:"total"`
	ExchangeRate                        decimal.Decimal `json:"exchangeRate"`
	Cancelled                           string          `json:"cancelled"`
	Discount                            decimal.Decimal `json:"discount"`
	WayOfPayment                        string          `json:"wayOfPayment"`
	PaymentMethod                       string          `json:"paymentMethod,omitempty"`
	CertificateNumber                   string          `json:"certificateNumber"`
	ConceptProductServiceKey            string          `json:"conceptProductServiceKey"`
	ConceptProductServiceKeyDescription string          `json:"conceptProductServiceKeyDescription"`
	Concepts                            string          `json:"concepts"`
	ConceptIdentificationNumber         string          `json:"conceptIdentificationNumber"`
	Unit                                string          `json:"unit"`
	ConceptQuantity                     decimal.Decimal `json:"conceptQuantity"`
	ConceptAmount                       decimal.Decimal `json:"conceptAmount"`
	ConceptUnitValue                    decimal.Decimal `json:"conceptUnitValue"`
	TransferredIva                      decimal.Decimal `json:"transferredIva"`
	TransferredIeps                     decimal.Decimal `json:"transferredIeps"`
	TransferredBase                     decimal.Decimal `json:"transferredBase"`
	TransferredTax                      decimal.Decimal `json:"transferredTax"`
	WithholdingTax                      decimal.Decimal `json:"withholdingTax"`
	WithholdingIsr                      decimal.Decimal `json:"withholdingIsr"`
	Base0Iva                            decimal.Decimal `json:"base0Iva"`
	Base8Iva                            decimal.Decimal `json:"base8Iva"`
	Base16Iva                           decimal.Decimal `json:"base16Iva"`
	BaseExemptIva                       decimal.Decimal `json:"baseExemptIva"`
	BaseIeps                            decimal.Decimal `json:"baseIeps"`
	IepsRateOrFee                       decimal.Decimal `json:"iepsRateOrFee"`
	Ieps                                decimal.Decimal `json:"ieps"`
	VatRateOrFee                        decimal.Decimal `json:"vatRateOrFee"`
	FileName                            string          `json:"fileName"`
	IsValid                             bool            `json:"isValid"`
}

// Validate implementa che
func (r *CFDIDetails2022) Validate() error {
	// Ejemplo de regla: UUID no vacío
	if r.Uuid == uuid.Nil {
		return fmt.Errorf("uuid vacío")
	}
	// Ejemplo: CfdiId no vacío
	if r.CfdiId == "" {
		return fmt.Errorf("cfdiId vacío")
	}
	// Aquí puedes añadir más reglas de validación...
	return nil
}

func main() {
	excelPath := flag.String("excel", "", "ruta al archivo Excel .xlsx")
	sheet := flag.String("sheet", "", "nombre de la hoja en el Excel")
	jsonDir := flag.String("jsondir", "", "directorio que contiene archivos JSON")
	flag.Parse()

	if *excelPath == "" || *sheet == "" || *jsonDir == "" {
		fmt.Println("Uso: cfdi-validator -excel input.xlsx -sheet CFDIDETAILS2022 -jsondir ./jsons")
		os.Exit(1)
	}

	excelCount, err := countExcelRowsStreaming(*excelPath, *sheet)
	if err != nil {
		log.Fatalf("Error leyendo Excel: %v", err)
	}
	fmt.Printf("Filas en Excel: %d\n", excelCount)

	jsonCount, err := countAllJSONObjects(*jsonDir)
	if err != nil {
		log.Fatalf("Error procesando JSONs: %v", err)
	}
	fmt.Printf("Objetos totales en JSONs: %d\n", jsonCount)

	if excelCount != jsonCount {
		log.Fatalf("❌ Mismatch: Excel=%d vs JSONs=%d", excelCount, jsonCount)
	}
	fmt.Println("✅ Coinciden filas Excel y total de objetos JSON.")
}

// countExcelRowsStreaming cuenta filas (sin cabecera) en modo streaming.
func countExcelRowsStreaming(path, sheetName string) (int, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	rows, err := f.Rows(sheetName)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// saltar cabecera
	if rows.Next() { /* descartar primera fila */
	}

	cnt := 0
	for rows.Next() {
		cnt++
	}
	return cnt, rows.Error()
}

// countAllJSONObjects recorre jsonDir y suma objetos de cada .json.
func countAllJSONObjects(jsonDir string) (int, error) {
	total := 0
	err := filepath.WalkDir(jsonDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		cnt, err := countJSONObjects(path)
		if err != nil {
			return fmt.Errorf("en %s: %w", d.Name(), err)
		}
		fmt.Printf(" → %s: %d objetos válidos\n", d.Name(), cnt)
		total += cnt
		return nil
	})
	return total, err
}

// countJSONObjects cuenta y valida elementos de un JSON (array o NDJSON).
func countJSONObjects(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	// saltar espacios en blanco iniciales
	for {
		b, err := r.Peek(1)
		if err != nil {
			return 0, err
		}
		if b[0] <= ' ' {
			_, _ = r.ReadByte()
			continue
		}
		break
	}

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	b, _ := r.Peek(1)
	cnt := 0

	switch b[0] {
	case '[':
		// JSON array
		if _, err := dec.Token(); err != nil {
			return 0, err
		}
		for dec.More() {
			var rec CFDIDetails2022
			if err := dec.Decode(&rec); err != nil {
				log.Printf("ERROR decode %s: registro %d: %v", path, cnt+1, err)
				continue
			}
			if err := rec.Validate(); err != nil {
				log.Printf("ERROR validación %s: registro %d: %v", path, cnt+1, err)
				continue
			}
			cnt++
		}
		if _, err := dec.Token(); err != nil {
			return cnt, err
		}

	case '{':
		// NDJSON (un objeto por línea)
		for {
			var rec CFDIDetails2022
			if err := dec.Decode(&rec); err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("ERROR decode %s: registro %d: %v", path, cnt+1, err)
				continue
			}
			if err := rec.Validate(); err != nil {
				log.Printf("ERROR validación %s: registro %d: %v", path, cnt+1, err)
				continue
			}
			cnt++
		}

	default:
		return 0, fmt.Errorf("archivo %s no parece ni array JSON ni NDJSON", path)
	}

	return cnt, nil
}
