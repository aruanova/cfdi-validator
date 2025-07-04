package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"os"
	"path/filepath"
)

// CFDIDetails2022 representa la estructura de cada objeto JSON.
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

// excelHeaderByJSONField mapea cada campo JSON a su encabezado en el Excel
var excelHeaderByJSONField = map[string]string{
	"cfdiId":                              "ID",
	"status":                              "ESTATUS",
	"cfdiRelationType":                    "TIPO RELACION",
	"relatedCfdi":                         "CFDI RELACIONADO",
	"type":                                "TIPO COMPROBANTE",
	"uuid":                                "UUID",
	"series":                              "SERIE",
	"reference":                           "FOLIO",
	"emitterRfc":                          "EMISOR",
	"emitterCompanyName":                  "RAZON SOCIAL EMISOR",
	"emitterPostalCode":                   "Emisor Codigo Postal",
	"cfdiUsage":                           "USO CFDI",
	"receiverRfc":                         "RECEPTOR",
	"receiverCompanyName":                 "RAZON SOCIAL RECEPTOR",
	"receptorPostalCode":                  "Receptor Codigo Postal",
	"origin":                              "ORIGEN",
	"currency":                            "MONEDA",
	"stampedDate":                         "FECHA TIMBRADO",
	"invoiceDate":                         "FECHA FACTURACION",
	"grouping":                            "AGRUPACION",
	"iva":                                 "IVA",
	"subTotal":                            "SUBTOTAL",
	"total":                               "TOTAL",
	"exchangeRate":                        "TIPO CAMBIO",
	"cancelled":                           "CANCELADO",
	"discount":                            "DESCUENTO",
	"wayOfPayment":                        "FORMA DE PAGO",
	"certificateNumber":                   "No CERTIFICADO",
	"conceptProductServiceKey":            "ClaveProducto o Servicio",
	"conceptProductServiceKeyDescription": "DescClaveProducto o Servicio",
	"concepts":                            "DESCRIPCION",
	"conceptIdentificationNumber":         "No IDENDIFICACION",
	"unit":                                "UNIDAD",
	"conceptQuantity":                     "CANTIDAD",
	"conceptAmount":                       "IMPORTE",
	"conceptUnitValue":                    "VALOR UNITARIO",
	"transferredIva":                      "IVA TRASLADO",
	"transferredIeps":                     "IEPS TRASLADO",
	"transferredBase":                     "TRASLADO BASE",
	"withholdingTax":                      "IMPUESTO RETENIDO",
	"transferredTax":                      "IMPUESTO TRASLADO",
	"base16Iva":                           "IVA BASE 16",
	"base8Iva":                            "IVA BASE 8",
	"base0Iva":                            "IVA BASE 0",
	"vatRateOrFee":                        "TASA O CUOTA IVA",
	"baseExemptIva":                       "BASE IVA Exento",
	"iepsRateOrFee":                       "TASA O CUOTA IEPS",
	"ieps":                                "IEPS",
	"baseIeps":                            "BASE IEPS",
	"withholdingIsr":                      "ISR RETENIDO",
	"paymentMethod":                       "METODO DE PAGO",
}

// Validate aplica reglas sencillas de validación sobre un registro.
func (r *CFDIDetails2022) Validate() error {
	if r.CfdiId == "" {
		return fmt.Errorf("cfdiId vacío")
	}
	if r.Uuid == uuid.Nil {
		return fmt.Errorf("uuid inválido")
	}
	// Agrega aquí más validaciones según sea necesario...
	return nil
}

// loadExcelKeys lee la columna keyName de Excel y devuelve un set de valores.
func loadExcelKeys(path, sheet, keyName string) (map[string]struct{}, error) {
	log.Printf("▶ Abriendo Excel: %s (hoja %s)", path, sheet)
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("error abriendo Excel: %w", err)
	}
	defer f.Close()

	rows, err := f.Rows(sheet)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo filas de la hoja %q: %w", sheet, err)
	}
	defer rows.Close()

	// 1) Leemos la primera fila como encabezado
	if !rows.Next() {
		return nil, fmt.Errorf("la hoja %q está vacía", sheet)
	}
	header, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error leyendo encabezado de %q: %w", sheet, err)
	}
	log.Printf("   Encabezados detectados: %v", header)

	// 2) Buscamos el índice de la columna keyName
	log.Printf("   Buscando columna %q en el encabezado...", keyName)
	keyIdx := -1
	for i, h := range header {
		if h == keyName {
			keyIdx = i
			break
		}
	}
	if keyIdx < 0 {
		return nil, fmt.Errorf("columna %q no encontrada entre %v", keyName, header)
	}
	log.Printf("   Columna %q encontrada en índice %d", keyName, keyIdx)

	// 3) Iteramos el resto de filas y vamos acumulando
	set := make(map[string]struct{}, 1_000_000)
	count := 0
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		if keyIdx < len(row) && row[keyIdx] != "" {
			set[row[keyIdx]] = struct{}{}
		}
		count++
		if count%200000 == 0 {
			log.Printf("   → %d filas leídas del Excel...", count)
		}
	}
	log.Printf("✔ Cargados %d CFDiIds únicos desde el Excel", len(set))
	return set, nil
}

// countExcelRowsStreaming cuenta las filas (sin cabecera) de la hoja Excel.
func countExcelRowsStreaming(path, sheet string) (int, error) {
	log.Printf("▶ Contando filas en Excel: %s (hoja %s)", path, sheet)
	f, err := excelize.OpenFile(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	rows, err := f.Rows(sheet)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// Saltar cabecera
	if rows.Next() {
		log.Printf("   Cabecera saltada")
		// nada
	}

	cnt := 0
	for rows.Next() {
		cnt++
		if cnt%250000 == 0 {
			log.Printf("   → %d filas contadas hasta ahora...", cnt)
		}
	}
	log.Printf("✔ Total de filas en Excel: %d", cnt)
	return cnt, rows.Error()
}

// countJSONObjects cuenta y valida objetos JSON en path, chequeando también
// los primeros 5 cfdiId contra excelKeys.
func countJSONObjects(path string, excelKeys map[string]struct{}) (int, error) {
	log.Printf("▶ Procesando JSON: %s", path)
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	// Saltar espacios en blanco iniciales
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

	// Detectar formato (array vs NDJSON)
	b, _ := r.Peek(1)
	var format string
	switch b[0] {
	case '[':
		format = "array"
	case '{':
		format = "ndjson"
	default:
		return 0, fmt.Errorf("formato JSON desconocido en %s (byte %#x)", path, b[0])
	}
	log.Printf("   Formato detectado: %s", format)

	cnt := 0
	switch format {
	case "array":
		// JSON array
		if _, err := dec.Token(); err != nil {
			return 0, err
		}
		for dec.More() {
			var rec CFDIDetails2022
			if err := dec.Decode(&rec); err != nil {
				log.Printf("   ERROR decode %s registro %d: %v", filepath.Base(path), cnt+1, err)
				continue
			}
			// validar los primeros 5 en Excel
			if cnt < 5 {
				if _, ok := excelKeys[rec.CfdiId]; !ok {
					log.Printf("   ⚠ cfdiId %q NO encontrado en Excel", rec.CfdiId)
				}
			}
			// validación de tipo
			if err := rec.Validate(); err != nil {
				log.Printf("   ERROR validación registro %d: %v", cnt+1, err)
			}
			cnt++
			if cnt%50000 == 0 {
				log.Printf("   → %d objetos leídos de %s...", cnt, filepath.Base(path))
			}
		}
		if _, err := dec.Token(); err != nil {
			return cnt, err
		}

	case "ndjson":
		// NDJSON (un objeto por línea)
		for {
			var rec CFDIDetails2022
			if err := dec.Decode(&rec); err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("   ERROR decode %s registro %d: %v", filepath.Base(path), cnt+1, err)
				continue
			}
			if cnt < 5 {
				if _, ok := excelKeys[rec.CfdiId]; !ok {
					log.Printf("   ⚠ cfdiId %q NO encontrado en Excel", rec.CfdiId)
				}
			}
			if err := rec.Validate(); err != nil {
				log.Printf("   ERROR validación registro %d: %v", cnt+1, err)
			}
			cnt++
			if cnt%50000 == 0 {
				log.Printf("   → %d objetos leídos de %s...", cnt, filepath.Base(path))
			}
		}
	}

	log.Printf("✔ Finalizado %s: %d objetos válidos", filepath.Base(path), cnt)
	return cnt, nil
}

// countAllJSONObjects recorre el directorio y suma los conteos de cada JSON.
func countAllJSONObjects(jsonDir string, excelKeys map[string]struct{}) (int, error) {
	log.Printf("▶ Barrido de JSONs en directorio: %s", jsonDir)
	total := 0
	err := filepath.WalkDir(jsonDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		cnt, err := countJSONObjects(path, excelKeys)
		if err != nil {
			return fmt.Errorf("en %s: %w", d.Name(), err)
		}
		log.Printf("   → %s: %d objetos", d.Name(), cnt)
		total += cnt
		return nil
	})
	log.Printf("✔ Total de objetos JSON leídos: %d", total)
	return total, err
}

func main() {
	excelPath := flag.String("excel", "", "ruta al archivo Excel .xlsx")
	sheet := flag.String("sheet", "", "nombre de la hoja")
	jsonDir := flag.String("jsondir", "", "directorio con JSONs")
	flag.Parse()

	if *excelPath == "" || *sheet == "" || *jsonDir == "" {
		fmt.Println("Uso: cfdi-validator -excel input.xlsx -sheet Hoja1 -jsondir ./jsons")
		os.Exit(1)
	}

	// 1) Leer cfdiId del Excel
	headerName := excelHeaderByJSONField["cfdiId"]
	excelKeys, err := loadExcelKeys(*excelPath, *sheet, headerName)
	if err != nil {
		log.Fatalf("Error cargando claves de Excel: %v", err)
	}

	// 2) Contar filas en Excel
	excelCount, err := countExcelRowsStreaming(*excelPath, *sheet)
	if err != nil {
		log.Fatalf("Error leyendo Excel: %v", err)
	}
	fmt.Printf("Filas en Excel: %d\n", excelCount)

	// 3) Contar y validar JSONs
	jsonCount, err := countAllJSONObjects(*jsonDir, excelKeys)
	if err != nil {
		log.Fatalf("Error procesando JSONs: %v", err)
	}
	fmt.Printf("Objetos totales en JSONs: %d\n", jsonCount)

	if excelCount != jsonCount {
		log.Fatalf("❌ Mismatch: Excel=%d vs JSONs=%d", excelCount, jsonCount)
	}
	fmt.Println("✅ Coinciden filas Excel y total de objetos JSON.")
}
