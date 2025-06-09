package commands

// logs flow.
// message added/popped to queue --> Add the same message into the genie for live log processing.
// After each pop from the live data, we will store the logs into a text file as a datastore.
// During each genie server start, we will remove the data completely.
