package likelihood

import "gonum.org/v1/gonum/mat"

// Statistics
type Statistics struct {
	Mean        *mat.VecDense
	Covariance  *mat.SymDense
	choleskyCov *mat.Cholesky
	inverseCov  *mat.Dense // shouldn't this also be Symmetric?
}

func (s *Statistics) SetCovariance(cov *mat.SymDense) {
	s.Covariance = cov
	ok := s.choleskyCov.Factorize(cov)
	if !ok {
		panic("couldn't set covariance")
	}
	err := s.inverseCov.Inverse(cov)
	if err != nil {
		panic(err)
	}
}

func (s *Statistics) GetCholeskyCovariance() *mat.Cholesky {
	return s.choleskyCov
}

func (s *Statistics) GetInverseCovariance() *mat.Dense {
	return s.inverseCov
}
