# Setting Up Automatic Focus Mode Control

Rune can automatically enable and disable Do Not Disturb (Focus mode) when you start and stop work sessions. To enable this functionality, you need to create two shortcuts in the macOS Shortcuts app.

## Setup Instructions

### 1. Open the Shortcuts App
- Press `Cmd + Space` and search for "Shortcuts"
- Or find it in your Applications folder

### 2. Create "Turn On Do Not Disturb" Shortcut
1. Click the "+" button to create a new shortcut
2. Search for "Set Focus" in the actions library
3. Drag "Set Focus" into your shortcut
4. In the Set Focus action, choose "Do Not Disturb" from the dropdown
5. Name your shortcut "Turn On Do Not Disturb" (exact name required)
6. Click "Done"

### 3. Create "Turn Off Do Not Disturb" Shortcut
1. Click the "+" button to create another new shortcut
2. Search for "Set Focus" in the actions library
3. Drag "Set Focus" into your shortcut
4. In the Set Focus action, choose "Turn Off Focus" from the dropdown
5. Name your shortcut "Turn Off Do Not Disturb" (exact name required)
6. Click "Done"

### 4. Test the Setup
Run `rune start` and `rune stop` to verify that Focus mode is automatically controlled.

## Alternative Setup

If you prefer different shortcut names or want to use a different Focus mode (like "Work" instead of "Do Not Disturb"), you can:

1. Create shortcuts with your preferred names
2. Update the shortcut names in the Rune source code in `internal/dnd/dnd.go`
3. Rebuild Rune with `go build`

## Troubleshooting

- **Shortcuts not working**: Make sure the shortcut names match exactly: "Turn On Do Not Disturb" and "Turn Off Do Not Disturb"
- **Permission issues**: The first time you run the shortcuts, macOS may ask for permissions
- **Focus modes not available**: Ensure you're running macOS 12 (Monterey) or later
- **Custom Focus modes**: You can modify the shortcuts to use custom Focus modes instead of the default "Do Not Disturb"

## Manual Control

If you prefer to control Focus mode manually, you can:
- Use Control Center: Click the Control Center icon → Focus → Do Not Disturb
- Use Siri: "Hey Siri, turn on Do Not Disturb"
- Set up a keyboard shortcut in System Settings → Keyboard → Keyboard Shortcuts → Mission Control