package react

import (
	"strings"
	"testing"
)

func TestAppendStepsFromText_Action(t *testing.T) {

	for _, tt := range []struct {
		name string
		text string
		want []Entry
	}{
		{
			name: "basic",
			text: `Question: Is the Python version used by the interpreter a stable release?
Assumption: I can use the Python interpreter
Thought: The version of the Python interpreter can be determined using the sys module in Python.

Action: python
import sys
sys.version`,
			want: []Entry{
				{Tag: Tags.Question, Content: "Is the Python version used by the interpreter a stable release?"},
				{Tag: Tags.Assumption, Content: "I can use the Python interpreter"},
				{Tag: Tags.Thought, Content: "The version of the Python interpreter can be determined using the sys module in Python."},
				{Tag: Tags.Action, Argument: "python", Content: "import sys\nsys.version"},
			},
		},
		{
			name: "line containing only tabs, empty lines",
			text: `

Question: Is the Python version used by the interpreter a stable release?
Thought: The version of the Python interpreter can be determined using the sys module in Python.


Action: python
import sys
sys.version`,
			want: []Entry{
				{Tag: Tags.Question, Content: "Is the Python version used by the interpreter a stable release?"},
				{Tag: Tags.Thought, Content: "The version of the Python interpreter can be determined using the sys module in Python."},
				{Tag: Tags.Action, Argument: "python", Content: "import sys\nsys.version"},
			},
		},
		{
			name: "thoughts, questions, actions",
			text: `Thought: The Python interpreter is using version 3.8.5.
Question: Is Python version 3.8.5 a stable release?
Thought: Stable releases of Python usually have a version number with two parts (major.minor) or three parts (major.minor.micro) if the micro version is zero. If the micro version is greater than zero, it is usually a bug fix release which is also considered stable.
Action: python
version_parts = tuple(map(int, '3.8.5'.split('.')))
len(version_parts) in {2, 3} and (len(version_parts) != 3 or version_parts[2] == 0)`,
			want: []Entry{
				{Tag: Tags.Thought, Content: "The Python interpreter is using version 3.8.5."},
				{Tag: Tags.Question, Content: "Is Python version 3.8.5 a stable release?"},
				{Tag: Tags.Thought, Content: "Stable releases of Python usually have a version number with two parts (major.minor) or three parts (major.minor.micro) if the micro version is zero. If the micro version is greater than zero, it is usually a bug fix release which is also considered stable."},
				{Tag: Tags.Action, Argument: "python", Content: "version_parts = tuple(map(int, '3.8.5'.split('.')))\nlen(version_parts) in {2, 3} and (len(version_parts) != 3 or version_parts[2] == 0)"},
			}},
		{
			name: "python indentation",
			text: `Question: Is the Python version used by the interpreter a stable release?
Thought: The version of the Python interpreter can be determined using the sys module in Python.

Action: python
def foo():
    return 1`,

			want: []Entry{
				{Tag: Tags.Question, Content: "Is the Python version used by the interpreter a stable release?"},
				{Tag: Tags.Thought, Content: "The version of the Python interpreter can be determined using the sys module in Python."},
				{Tag: Tags.Action, Argument: "python", Content: "def foo():\n    return 1"},
			},
		},

		{
			name: "python don't lose empty lines",
			text: `Question: Is the Python version used by the interpreter a stable release?
Thought: The version of the Python interpreter can be determined using the sys module in Python.

Action: python
def foo():
    return 1

foo()`,

			want: []Entry{
				{Tag: Tags.Question, Content: "Is the Python version used by the interpreter a stable release?"},
				{Tag: Tags.Thought, Content: "The version of the Python interpreter can be determined using the sys module in Python."},
				{Tag: Tags.Action, Argument: "python", Content: "def foo():\n    return 1\n\nfoo()"},
			},
		},
		{
			name: "multiline question and thought",
			text: `Question: Is the Python version used by the interpreter a stable release?

A lot depends on that.
Thought: The version of the Python interpreter can be determined using the sys module in Python.

Let's give it a try.

Action: python
import sys
sys.version`,
			want: []Entry{
				{Tag: Tags.Question, Content: "Is the Python version used by the interpreter a stable release?\n\nA lot depends on that."},
				{Tag: Tags.Thought, Content: "The version of the Python interpreter can be determined using the sys module in Python.\n\nLet's give it a try."},
				{Tag: Tags.Action, Argument: "python", Content: "import sys\nsys.version"},
			},
		},
	} {

		steps, err := Parse(tt.text)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tt.name, err)
		}

		for _, step := range steps {
			step.Content = strings.TrimSpace(step.Content)
		}

		for i, step := range steps {
			if step.Content != strings.TrimSpace(step.Content) {
				t.Fatalf("%s: unexpected whitespace at step %d: %q", tt.name, i, step.Content)
			}
		}

		if len(steps) != len(tt.want) {
			t.Errorf("expected %d steps, got %d", len(tt.want), len(steps))
		}

		for i, step := range steps {
			if step.Tag != tt.want[i].Tag {
				t.Errorf("%s: unexpected step at index %d: type %s (want type %s)", tt.name, i, step.Tag, tt.want[i].Tag)
			}

			if step.Argument != tt.want[i].Argument {
				t.Errorf("%s: unexpected step at index %d: argument %q (want argument %q)", tt.name, i, step.Argument, tt.want[i].Argument)
			}

			if step.Content != tt.want[i].Content {
				t.Errorf("%s: unexpected step at index %d: content %q (want content %q)", tt.name, i, step.Content, tt.want[i].Content)
			}

		}
	}
}
