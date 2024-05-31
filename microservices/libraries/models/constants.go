package models

const Id = "IdField"
const Timestamp = "Timestamp"
const LastDestinationId = "LastDestinationId"
const LastDestinationTimestamp = "LastDestinationTimestamp"
const FullWithId = "FullWithId"
const WriteOnly = "WriteOnly"

// SAVE MODES
const Insert = "Insert"
const Upsert = "Upsert"

// ERROR MODES
const Ignore = "Ignore"
const StopOnError = "StopOnError"

// RUN MODES
const Static = "Static"
const StaticFilePath = "StaticFilePath"
const DisableMachineId = "DisableMachineId"

const MaxBatchSizeDefault = 5000

// PATHS
const OffsetsPath = "/offsets/"
const AssignmentsPath = "/assignments/"
const ErrorsPath = "/errors/"
