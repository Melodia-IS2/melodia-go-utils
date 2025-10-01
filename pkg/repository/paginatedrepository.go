package repository

import (
	"context"
	"math"
	"sort"
	"sync"

	"github.com/Melodia-IS2/melodia-go-utils/pkg/pages"
)

type fetchFn[T any] func(ctx context.Context, pagination pages.Pagination) ([]T, pages.PaginationResult, error)

type requestInfo struct {
	pagination pages.Pagination
	pageIndex  int
}

// pageResult contiene el resultado de un request específico
type pageResult[T any] struct {
	data             []T
	pageIndex        int
	paginationResult pages.PaginationResult
	err              error
}

func FetchPaginatedExternal[T any](ctx context.Context, fetchFn fetchFn[T], pagination pages.Pagination, externalPageSize uint) ([]T, pages.PaginationResult, error) {
	// 1. Determinar el tamaño de página efectivo a usar
	effectivePageSize := pagination.PageSize
	if pagination.PageSize > externalPageSize {
		effectivePageSize = externalPageSize
	}

	// 2. Hacer el primer request para obtener información total
	firstPageRequest := pages.Pagination{
		Page:     pagination.Page,
		PageSize: effectivePageSize,
	}

	firstPageResult, paginationResult, err := fetchFn(ctx, firstPageRequest)
	if err != nil {
		// Ante cualquier respuesta distinta a 200, retornarla
		return nil, pages.PaginationResult{}, err
	}

	// Si el PageSize interno es menor o igual que el externo, retornar directamente
	if pagination.PageSize <= externalPageSize {
		return firstPageResult, paginationResult, nil
	}

	// 3. Precalcular todos los requests necesarios
	recordsNeeded := pagination.PageSize
	totalAvailable := paginationResult.TotalRecords
	startPage := pagination.Page

	// Calcular cuántos registros necesitamos en total (considerando registros disponibles)
	actualRecordsToFetch := recordsNeeded
	if totalAvailable < recordsNeeded {
		actualRecordsToFetch = totalAvailable
	}

	// Si ya tenemos todo lo que necesitamos con la primera página
	if uint(len(firstPageResult)) >= actualRecordsToFetch {
		// Pre-allocar y copiar solo los registros necesarios
		finalResult := make([]T, actualRecordsToFetch)
		if uint(len(firstPageResult)) > actualRecordsToFetch {
			copy(finalResult, firstPageResult[:actualRecordsToFetch])
		} else {
			copy(finalResult, firstPageResult)
		}

		finalPaginationResult := pages.PaginationResult{
			Page:         pagination.Page,
			PageSize:     pagination.PageSize,
			TotalRecords: paginationResult.TotalRecords,
			TotalPages:   uint(math.Ceil(float64(paginationResult.TotalRecords) / float64(pagination.PageSize))),
		}

		return finalResult, finalPaginationResult, nil
	}

	// 4. Calcular requests adicionales necesarios
	recordsAlreadyFetched := uint(len(firstPageResult))
	recordsRemaining := actualRecordsToFetch - recordsAlreadyFetched

	// Pre-calcular el número de requests adicionales necesarios
	numAdditionalRequests := int(math.Ceil(float64(recordsRemaining) / float64(effectivePageSize)))

	// Pre-allocar el slice con el tamaño exacto
	additionalRequests := make([]requestInfo, numAdditionalRequests)

	currentPage := startPage + 1
	pageIndex := 1 // El primer request ya se hizo, empezamos en índice 1
	requestIndex := 0

	for recordsRemaining > 0 && requestIndex < numAdditionalRequests {
		nextPageSize := effectivePageSize
		if recordsRemaining < effectivePageSize {
			nextPageSize = recordsRemaining
		}

		additionalRequests[requestIndex] = requestInfo{
			pagination: pages.Pagination{
				Page:     currentPage,
				PageSize: nextPageSize,
			},
			pageIndex: pageIndex,
		}

		recordsRemaining -= nextPageSize
		currentPage++
		pageIndex++
		requestIndex++
	}

	// 5. Ejecutar requests adicionales de forma concurrente
	var wg sync.WaitGroup
	resultsChan := make(chan pageResult[T], len(additionalRequests))

	for _, req := range additionalRequests {
		wg.Add(1)
		go func(request requestInfo) {
			defer wg.Done()

			data, pagResult, err := fetchFn(ctx, request.pagination)
			resultsChan <- pageResult[T]{
				data:             data,
				pageIndex:        request.pageIndex,
				paginationResult: pagResult,
				err:              err,
			}
		}(req)
	}

	// Cerrar el canal cuando todas las goroutines terminen
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 6. Recopilar resultados y verificar errores
	allResults := make([]pageResult[T], 0, len(additionalRequests))
	for result := range resultsChan {
		if result.err != nil {
			// Ante cualquier error, retornarlo inmediatamente
			return nil, pages.PaginationResult{}, result.err
		}
		allResults = append(allResults, result)
	}

	// 7. Ordenar resultados por pageIndex para mantener el orden correcto
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].pageIndex < allResults[j].pageIndex
	})

	// 8. Combinar todos los resultados
	// Pre-calcular tamaño total necesario
	totalSize := len(firstPageResult)
	for _, result := range allResults {
		totalSize += len(result.data)
	}

	// Pre-allocar slice final con tamaño exacto
	finalResult := make([]T, totalSize)

	// Copiar primera página
	offset := copy(finalResult, firstPageResult)

	// Copiar resultados adicionales usando copy para máximo rendimiento
	for _, result := range allResults {
		offset += copy(finalResult[offset:], result.data)
	}

	// 9. Construir el resultado final de paginación
	finalPaginationResult := pages.PaginationResult{
		Page:         pagination.Page,
		PageSize:     pagination.PageSize,
		TotalRecords: paginationResult.TotalRecords,
		TotalPages:   uint(math.Ceil(float64(paginationResult.TotalRecords) / float64(pagination.PageSize))),
	}

	return finalResult, finalPaginationResult, nil
}
