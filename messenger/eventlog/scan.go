package eventlog

// Scan helps scanning an event log after a given version.
// Scan will read the event stream in batches the size of the given buffer,
// call onEvent for every read event and stop if onEvent returns false
// returning the version at which the scan stopped.
// Scan will allocate a new buffer if the length of the given buffer is <1
func Scan(
	log *EventLog,
	after Version,
	buffer []Event,
	onEvent func(Event) bool,
) (Version, error) {
	if len(buffer) < 1 {
		buffer = make([]Event, ScanDefaultBufferSize)
	}
	for {
		read, version, err := log.Read(after, buffer)
		if err != nil {
			return after, err
		}
		for i, e := range buffer[:read] {
			if !onEvent(e) {
				return Version(int(after) + i + 1), nil
			}
		}
		after += Version(read)
		if after == version {
			return version, nil
		}
	}
}

// ScanDefaultBufferSize defines the size of the buffer
// eventlog.Scan allocates when the provided buffer is nil
const ScanDefaultBufferSize = 64
