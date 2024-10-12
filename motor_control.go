package main

// Set the motor into driving forward
func driveForward(fullSpeed bool) {
	MIN1.High()
	MIN2.Low()
}

// Set the motor into driving backwards
func driveBackward(fullSpeed bool) {
	MIN1.Low()
	MIN2.High()
}

// Set the motor into full stop
func fullStop() {
	MIN1.Low()
	MIN2.Low()
}
