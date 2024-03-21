package main

import (
	"fmt"
)

// tip de data pentru un punct in plan
type Point struct {
	X, Y float64
}

// tip de data pentru un segment al splinei cubice
// contine coeficientii si valoarea x de la inceputul segmentului
type SplineSegment struct {
	a, b, c, d, x float64
}

// colectie de segmente consecutive (practic tot spline-ul cubic )
type CubicSpline struct {
	Segments []SplineSegment
}

// Calculeaza spline-ul cubic pentru un set de puncte in plan
// Returneaza o structura CubicSpline
func NaturalCubicSplines(points []Point) CubicSpline {
	n := len(points) - 1
	h := make([]float64, n) //distanta intre noduri consecutive

	//Calculul lui h
	for i := 0; i < n; i++ {
		h[i] = points[i+1].X - points[i].X
	}

	// Calculul lui alpha
	alpha := make([]float64, n)
	for i := 1; i < n; i++ {
		alpha[i] = 3/h[i]*(points[i+1].Y-points[i].Y) - 3/h[i-1]*(points[i].Y-points[i-1].Y)
	}

	//vectori ajutatori
	l := make([]float64, n+1)
	mu := make([]float64, n)
	z := make([]float64, n+1)
	c := make([]float64, n+1)

	// Setarile initiale
	l[0] = 1
	z[0] = 0
	c[n] = 0

	//Calculul Coeficientilor c pentru spline
	for i := 1; i < n; i++ {
		l[i] = 2*(points[i+1].X-points[i-1].X) - h[i-1]*mu[i-1]
		mu[i] = h[i] / l[i]
		z[i] = (alpha[i] - h[i-1]*z[i-1]) / l[i]
	}

	l[n] = 1
	z[n] = 0

	//init Coeficientii a,b,d (calculati pe baza lui c)
	b := make([]float64, n)
	d := make([]float64, n)
	a := make([]float64, n+1)

	//calculul final
	for j := n - 1; j >= 0; j-- {
		c[j] = z[j] - mu[j]*c[j+1]
		b[j] = (points[j+1].Y-points[j].Y)/h[j] - h[j]*(c[j+1]+2*c[j])/3
		d[j] = (c[j+1] - c[j]) / (3 * h[j])
		a[j] = points[j].Y
	}

	// Crearea spline-urilor
	splines := make([]SplineSegment, n)
	for i := 0; i < n; i++ {
		splines[i] = SplineSegment{a: a[i], b: b[i], c: c[i], d: d[i], x: points[i].X}
	}

	return CubicSpline{Segments: splines}
}

// Calculeaza valoarea spline-ului intr-un punct x dat (in cazul functiei main este z)
func (cs *CubicSpline) Evaluate(x float64) float64 {
	//se cauta unde se incadreaza x
	for i, segment := range cs.Segments {
		//verificam daca x se afla in intervalul [segment.x, cs.Segments[i+1].x) pentru segmentele intermediare
		//sau daca e mai mare sau egal cu segment.x pentru ultimul segment
		if (i < len(cs.Segments)-1 && x >= segment.x && x < cs.Segments[i+1].x) || (i == len(cs.Segments)-1 && x >= segment.x) {
			//dif dintre x si capatul stang
			dx := x - segment.x
			//se aplica formula polinomului
			return segment.a + segment.b*dx + segment.c*dx*dx + segment.d*dx*dx*dx
		}
	}
	return 0
}

func main() {
	//punctele in plan
	points := []Point{{1, 2}, {2, 3}, {3, 5}, {4, 7}, {5, 11}}
	//valorile z pentru care se calculeaza functiile
	z := []float64{1.5, 2.5, 3.5, 4.5}

	//lista de segment
	cubicSpline := NaturalCubicSplines(points)

	//calculul functiilor pe valorile Z
	for _, zi := range z {
		yi := cubicSpline.Evaluate(zi)
		fmt.Printf("Spline(%.2f) = %.2f\n", zi, yi)
	}
}
