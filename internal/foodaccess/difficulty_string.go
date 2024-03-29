// Code generated by "stringer -type=Difficulty -linecomment"; DO NOT EDIT.

package foodaccess

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DifficultyEasy-0]
	_ = x[DifficultyAverage-1]
	_ = x[DifficultyHard-2]
}

const _Difficulty_name = "facilemoyendifficile"

var _Difficulty_index = [...]uint8{0, 6, 11, 20}

func (i Difficulty) String() string {
	if i < 0 || i >= Difficulty(len(_Difficulty_index)-1) {
		return "Difficulty(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Difficulty_name[_Difficulty_index[i]:_Difficulty_index[i+1]]
}
