on run {partialNumberSought}
  set listOfFolks to {}
  tell application "Address Book" to set {phoneNumbers, peoplesNames} to {value of phones, name} of people
  repeat with p from 1 to (count peoplesNames)
     repeat with aPhone in item p of phoneNumbers
         set aPhone to do shell script "echo " & quoted form of aPhone & " | tr -dc '[:alnum:]' | tr '[:upper:]' '[:lower:]'"
         if (aPhone contains partialNumberSought) then
             set end of listOfFolks to item p of peoplesNames
             -- set output to ( item p of peoplesNames )
             -- log output
             -- set output to "TEST"
             -- do shell script "echo " & quoted form of output
             exit repeat
         end if
     end repeat
  end repeat
return listOfFolks
end run