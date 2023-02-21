package main

type KLineInterval int

const (
	KLineInterval1s = iota
	KLineInterval1m
	KLineInterval3m
	KLineInterval5m
	KLineInterval15m
	KLineInterval30m
	KLineInterval1h
	KLineInterval2h
	KLineInterval4h
	KLineInterval6h
	KLineInterval8h
	KLineInterval12h
	KLineInterval1d
	KLineInterval3d
	KLineInterval1w
	KLineInterval1M
)

func (k KLineInterval) String() string {
	return [...]string{
		"1s",
		"1m",
		"3m",
		"5m",
		"15m",
		"30m",
		"1h",
		"2h",
		"4h",
		"6h",
		"8h",
		"12h",
		"1d",
		"3d",
		"1w",
		"1M",
	}[k]
}
