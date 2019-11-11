package weathercomparer

//Temperature represents a float64 to contain a temperature value
type Temperature float64

func (temp Temperature) toCelsius() float64 {
	return (float64(temp) - 32) * 5 / 9
}