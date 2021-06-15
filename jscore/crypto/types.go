package js

type ECDSANamedCurve int

const (
	ECDSAP256 ECDSANamedCurve = iota
	ECDSAP384 ECDSANamedCurve = iota
	ECDSAP521 ECDSANamedCurve = iota
)

type ECDSAKeyUsage int

const (
	ECDSASign   ECDSAKeyUsage = iota
	ECDSAVerify ECDSAKeyUsage = iota
)
