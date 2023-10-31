package hooker.tracker.slack

#Example of handling tracker event

title:=sprintf("Tracker Detection - %s", [input.SigMetadata.Name])

result:= res {
 res:= [
 	{ "type":"section",
	  "text": {"type":"mrkdwn","text": sprintf("*Rule Description:* %s", [input.SigMetadata.Description])}},
 	{ "type":"section",
	  "text": {"type":"mrkdwn","text": sprintf("*Detection:* %s", [input.Context.processName])}},
 	{ "type":"section",
	  "text": {"type":"mrkdwn","text": sprintf("*MITRE Details:* %v", [input.SigMetadata.Properties])}},
	{ "type":"section",
	  "text": {"type":"mrkdwn","text": sprintf("*Severity:* %v", [input.SigMetadata.Properties.Severity])}}
 ]
}