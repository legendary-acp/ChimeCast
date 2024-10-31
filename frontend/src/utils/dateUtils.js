export default function formatTimestamp(timestampStr) {
    // Create a Date object from the timestamp string
    const date = new Date(timestampStr);

    // Options for formatting the date
    const options = {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: 'numeric',
        minute: 'numeric',
        hour12: true
    };

    // Format the date to a human-readable string
    return date.toLocaleString('en-US', options);
}