# pd-oncall-to-slack-channel

This is a simple slack bot with a single purpose: setting a slack channel's topic
to the name of engineer on call.

The mapping is between a pagerduty schedule and a slack channel, so you'll the id
for both.

**PS**: If the channel is private, you'll need to change slack endpoints to `groups` 
instead of `channels`.

### Requirements
- A Slack token set as `SLACK_TOKEN` environment variable
- A PagerDuty token set as `PD_TOKEN` environment variable
- A slack channel ID and corresponding PagerDuty schedule ID (hardcoded in the code for now)

### Todo
- [ ] Make it work with more than one channel<->schedule mapping
- [ ] Read channel<->schedule mappings from a file
- [ ] Break the one-to-one mapping, make it possible to map a schedule to multiple channels
- [ ] Better error handling
- [ ] Don't require slack channel Id, use channel name instead
- [ ] Work with both private and public channels
