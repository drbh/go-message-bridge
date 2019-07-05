on run {targetBuddyPhone, imageFilePath }
    tell application "Messages"
    set targetService to 1st service whose service type = iMessage
        set filen to (imageFilePath as POSIX file)
        set theAttachment to filen
        
        set targetBuddy to buddy targetBuddyPhone of targetService
        send theAttachment to targetBuddy
    end tell
end run