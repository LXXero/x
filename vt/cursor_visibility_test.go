package vt

import "testing"

// TestRestoreCursorFiresVisibility guards the DECRC visibility bug: a
// save -> hide (?25l) -> restore (DECRC) -> show (?25h) sequence — as
// Powerlevel10k's instant prompt emits — must leave an embedder that
// mirrors the CursorVisibility callback showing the cursor. Before the
// fix, RestoreCursor reset Hidden silently, so the trailing ?25h saw no
// change and never re-fired visible=true, stranding the mirror hidden.
func TestRestoreCursorFiresVisibility(t *testing.T) {
	term := newTestTerminal(t, 20, 3)

	var events []bool
	term.SetCallbacks(Callbacks{
		CursorVisibility: func(visible bool) { events = append(events, visible) },
	})

	// DECSC (ESC 7) save, hide, DECRC (ESC 8) restore, show.
	term.WriteString("\x1b7")     // save cursor (visible)
	term.WriteString("\x1b[?25l") // hide
	term.WriteString("\x1b8")     // restore -> visible again
	term.WriteString("\x1b[?25h") // show

	if len(events) == 0 {
		t.Fatalf("no CursorVisibility callbacks fired")
	}
	if last := events[len(events)-1]; !last {
		t.Fatalf("final cursor visibility = %v, want true; events=%v", last, events)
	}
}

// TestHideThenShowFiresBothVisibilities is the plain (no save/restore)
// baseline: ?25l then ?25h must fire false then true.
func TestHideThenShowFiresBothVisibilities(t *testing.T) {
	term := newTestTerminal(t, 20, 3)
	var events []bool
	term.SetCallbacks(Callbacks{
		CursorVisibility: func(visible bool) { events = append(events, visible) },
	})
	term.WriteString("\x1b[?25l")
	term.WriteString("\x1b[?25h")
	if len(events) != 2 || events[0] != false || events[1] != true {
		t.Fatalf("events = %v, want [false true]", events)
	}
}
