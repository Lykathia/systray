//go:build (linux || freebsd || openbsd || netbsd) && !android

package systray

import "testing"

// TestResetMenu verifies that ResetMenu drops every tracked item, closes each
// item's ClickedCh, and clears the underlying layout tree.
func TestResetMenu(t *testing.T) {
	// menuItems is a package global shared across the suite; start from a clean slate.
	menuItemsLock.Lock()
	menuItems = make(map[uint32]*MenuItem)
	menuItemsLock.Unlock()

	parent := AddMenuItem("parent", "")
	child := parent.AddSubMenuItem("child", "")
	sibling := AddMenuItemCheckbox("sibling", "", true)
	items := []*MenuItem{parent, child, sibling}

	menuItemsLock.RLock()
	tracked := len(menuItems)
	menuItemsLock.RUnlock()
	if tracked != len(items) {
		t.Fatalf("tracked items before reset = %d, want %d", tracked, len(items))
	}

	ResetMenu()

	menuItemsLock.RLock()
	remaining := len(menuItems)
	menuItemsLock.RUnlock()
	if remaining != 0 {
		t.Errorf("tracked items after reset = %d, want 0", remaining)
	}

	for _, item := range items {
		select {
		case _, ok := <-item.ClickedCh:
			if ok {
				t.Errorf("ClickedCh for item %d should be closed", item.id)
			}
		default:
			t.Errorf("ClickedCh for item %d not closed", item.id)
		}
	}

	instance.menuLock.Lock()
	layoutChildren := len(instance.menu.V2)
	instance.menuLock.Unlock()
	if layoutChildren != 0 {
		t.Errorf("layout tree after reset has %d children, want 0", layoutChildren)
	}
}
