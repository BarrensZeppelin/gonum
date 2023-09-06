// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/lapack"
)

type Dlanhser interface {
	Dlanhs(norm lapack.MatrixNorm, n int, a []float64, lda int, work []float64) float64
}

func DlanhsTest(t *testing.T, impl Dlanhser) {
	const tol = 1e-15
	allwork := make([]float64, 9)
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{1, 2, 4, 9} {
		for _, lda := range []int{n, n + 5} {
			a := randomHessenberg(n, lda, rnd)
			for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxRowSum, lapack.MaxColumnSum, lapack.Frobenius} {
				var work []float64
				if norm == lapack.MaxColumnSum {
					work = allwork[:n]
				}
				for i := range work {
					work[i] = math.NaN()
				}
				want := dlange(norm, a.Rows, a.Cols, a.Data, a.Stride)
				got := impl.Dlanhs(norm, a.Rows, a.Data, lda, work)
				if math.Abs(want-got)/want > tol {
					t.Errorf("Case n=%v,lda=%v,norm=%v: unexpected result. Want %v, got %v.", n, lda, normToString(norm), want, got)
				}
			}
		}
	}
}