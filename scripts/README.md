# Scripts Directory

This directory contains utility scripts for testing, debugging, and development.

## Development Scripts

### watch.sh

Hot reload script that automatically rebuilds and restarts the server when Go files change.

**Usage:**
```bash
# Via Makefile (recommended)
make dev

# Direct usage
./scripts/watch.sh

# With arguments
./scripts/watch.sh --use_oauth2
```

**Features:**
- Monitors all `.go` files for changes
- Automatically rebuilds on file save
- Restarts server with new changes
- Shows build errors in console
- Clean shutdown with Ctrl+C

**How it works:**
1. Builds initial binary
2. Starts server in background
3. Watches for file modifications (every 1 second)
4. On change: kills old server, rebuilds, starts new server
5. Continues until you press Ctrl+C

## Test Scripts

### Phase 3 Test Scripts

Scripts created during Phase 3 development to test assignment management features:

- `phase3-task1.sh` through `phase3-task10.sh` - Individual task testing
- `test_phase3_complete.sh` - Complete Phase 3 functionality test

### Feature Test Scripts

- `test_instructor_assign_feature.sh` - Test instructor assignment workflow
- `test_student_dashboard.sh` - Test student dashboard functionality
- `test_complete_assign_feature.sh` - Complete assignment feature test
- `test_task5_progress_tracking.sh` - Progress tracking tests
- `test_view_progress.sh` - View progress feature test
- `verify_task5_integration.sh` - Integration verification

### Dashboard Test Scripts

- `complete_dashboard_test.sh` - Complete dashboard testing
- `student_dashboard_verification.sh` - Student dashboard verification
- `demo_assignment_flow.sh` - Demo of assignment workflow

### Debug Scripts

- `debug_assignment_detail.sh` - Debug assignment detail views
- `debug_isolation.sh` - Isolated debugging
- `debug_test.sh` - General debugging
- `debug_test8.sh` - Specific debug scenario

## Usage Notes

### Making Scripts Executable

If a script isn't executable, run:
```bash
chmod +x scripts/script-name.sh
```

### Running Test Scripts

Most test scripts expect the server to be running:
```bash
# Terminal 1: Start server
make run

# Terminal 2: Run test script
./scripts/test_phase3_complete.sh
```

### Script Dependencies

Some scripts require:
- `curl` - For HTTP requests (standard on macOS/Linux)
- `jq` - For JSON parsing (install with `brew install jq`)
- Server running on `localhost:8080`

## Creating New Scripts

When creating new scripts:

1. **Use descriptive names**: `test_feature_name.sh` or `debug_issue_name.sh`
2. **Add shebang**: Start with `#!/bin/bash`
3. **Make executable**: `chmod +x scripts/your-script.sh`
4. **Document usage**: Add comments explaining what the script does
5. **Handle cleanup**: Use traps for Ctrl+C cleanup if needed

Example:
```bash
#!/bin/bash
# Description: Test the new feature X
# Usage: ./scripts/test_feature_x.sh

# Cleanup on exit
trap 'echo "Cleaning up..."; exit' INT TERM

# Your script logic here
echo "Testing feature X..."
```

## Best Practices

1. **Don't commit temporary test data** - Use `.gitignore` for cookies, temp files
2. **Use relative paths** - Scripts should work from project root
3. **Check for running server** - Test scripts should verify server is accessible
4. **Print clear output** - Use echo with status indicators (✅ ❌ ⚠️)
5. **Exit on errors** - Use `set -e` or check command results

## Maintenance

### Outdated Scripts

If a script is no longer relevant:
1. Move to `scripts/archive/` directory (create if needed)
2. Or delete if the feature was removed
3. Update this README

### Script Organization

Current organization:
- Development: `watch.sh`
- Phase tests: `phase3-task*.sh`
- Feature tests: `test_*.sh`
- Debugging: `debug_*.sh`
- Demos: `demo_*.sh`
- Verification: `verify_*.sh`

## Troubleshooting

### "Permission denied"
```bash
chmod +x scripts/script-name.sh
```

### "Command not found: jq"
```bash
brew install jq
```

### "Connection refused"
Make sure the server is running:
```bash
make run
```

### Scripts hanging
Press `Ctrl+C` to stop. If server process remains:
```bash
pkill -f zipcodereader
```

## Resources

- [Bash Scripting Guide](https://www.gnu.org/software/bash/manual/)
- [curl Documentation](https://curl.se/docs/)
- [jq Manual](https://stedolan.github.io/jq/manual/)

---

**Last Updated**: January 2026