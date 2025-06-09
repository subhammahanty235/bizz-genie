package commands

// It will show the connected Packages, And the packages heartbeat status.
// We will identify the packages/instances with an instance id, assigned during server start.
// each instance that is inactive, will be removed after some failed heartbeats of 10times,
