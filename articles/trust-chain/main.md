# Trust Chain
This is a Service Composition pattern in which services aren't aware of each-other when communicating. They dispatch
tasks by broadcasting the request, other service(s) receive it and if the task is interesting they declare themselves as
being in-charge immediately. In fact task owner knows how many responses it has to collect before sending the final
response.

## How to achieve it
Suppose having a **task dispatcher** which broadcasts the task without having information about which unit(s)
are going to work on the task and some **worker**s which may work on some tasks. Task dispatch and gathering
the result(s) is asynchronous. **task dispatcher** and **worker**s share an array per each task accessible by
{TaskName}-{RequestID}. TaskName is equivalent to message topic in messaging systems and
RequestID is a uniquely generated id by **task dispatcher**.

Whenever **worker**s receive a task, in case they realized that the task is interesting for them which means
they want to process it, they _push_ to the shared array {TaskName}-{RequestID}. Value doesn't matter.

When each **worker** finishes working on the task it publishes its individual result so that the **task dispatcher**
receives it as one partial response while there may be other possible responses too.
In this step the **task dispatcher** *pop*s from the shared array {TaskName}-{RequestID}. If there were still remaining
items on the shared array it means there are **worker**s which are still working on the task
in fact the **task dispatcher** don't return the final response and just keeps this partial response somewhere to join
with other partial responses later. If there were no items left on the shared array it means
this was the partial response from last **worker** which means all workers have finished their jobs, and it's time to
prepare the final response by combining partial responses from all **worker**s. Regardless of combining logic,
final response is ready to be sent.

## Implement using Messaging & Redis List
This patterns let us handle synchronous communication (frontend-backend) by variable number of
asynchronous communications (service-service) which may differ due to the request payload.
In the diagram below communication between internal services are event-driven.
Any async messaging platform can play role of the Message BUS.
[![](https://mermaid.ink/img/pako:eNq1VMGK2zAQ_ZVBh9LS5GCbvfgQyGZp2FJDsOjNF1meJGJtyZXklLDsv3dkJ9nEm92UQi_CmvckvXlj3jOTpkKWMoe_OtQSH5TYWNEUGkB03uiuKdEWut9LbywsQDhY1Aq1D8VWWK-kaoX2MF8GbL56hKXw-Fvsx4Q84DlWyo2RLCAZOic2CPc_-RjnUSBwtDslEaI3cHwOx2_g5BxOhnYsSg92U36O7-4mcFy-BGgxnc2-zpcp3NdGPim9IdFkj-s7ni8JzVJYdWWt3Pao-uSR2lHzkIVdNiUqj1LgXemkVSVWFyQevbLi91nxKysZsyo8fzJUeETEPIVH7dGSZqzgE4klpd6A3yL8UEMjp4N5fyz-q2MnKjnwTWkygIjfTTnSMmgeuzGbBVNzbOs9rRLV7loTB4tJy8q0DtbWNGMB0SDgQblWeLlFBwJcV3rhnt6bwg3nzv0-_CA8-Tcjk9vuJNfdCb_KbXeOA77mzgU7v3TrI0HR_xrXNUE9f7qggaCu6GbXGu0wILRnE9agbYSqKJWeQ7FgdF2DBUvpc23CKApW6Bdihojiey1Z6m2HE9a1Fb1zyDCWrkXtqEqBQ8mVDUHX593LH6K5okU?type=png)](https://mermaid.live/edit#pako:eNq1VMGK2zAQ_ZVBh9LS5GCbvfgQyGZp2FJDsOjNF1meJGJtyZXklLDsv3dkJ9nEm92UQi_CmvckvXlj3jOTpkKWMoe_OtQSH5TYWNEUGkB03uiuKdEWut9LbywsQDhY1Aq1D8VWWK-kaoX2MF8GbL56hKXw-Fvsx4Q84DlWyo2RLCAZOic2CPc_-RjnUSBwtDslEaI3cHwOx2_g5BxOhnYsSg92U36O7-4mcFy-BGgxnc2-zpcp3NdGPim9IdFkj-s7ni8JzVJYdWWt3Pao-uSR2lHzkIVdNiUqj1LgXemkVSVWFyQevbLi91nxKysZsyo8fzJUeETEPIVH7dGSZqzgE4klpd6A3yL8UEMjp4N5fyz-q2MnKjnwTWkygIjfTTnSMmgeuzGbBVNzbOs9rRLV7loTB4tJy8q0DtbWNGMB0SDgQblWeLlFBwJcV3rhnt6bwg3nzv0-_CA8-Tcjk9vuJNfdCb_KbXeOA77mzgU7v3TrI0HR_xrXNUE9f7qggaCu6GbXGu0wILRnE9agbYSqKJWeQ7FgdF2DBUvpc23CKApW6Bdihojiey1Z6m2HE9a1Fb1zyDCWrkXtqEqBQ8mVDUHX593LH6K5okU)

## Prevent waiting forever
Check length of the shared array {TaskName}-{RequestID} a bit (500ms~1s) after task got dispatched. If it's greater than
zero it means there are at least one worker which is working on the task and **task dispatcher** has to wait for it.
If it's zero it means no one interested in the task and the **task dispatcher** has to return an empty response.

## Nesting
Not a big deal. Each **worker** can be a **task dispatcher** too for a 2nd task which is necessary to get done
in order to finish the 1st task.

## Compensation

## An example

## Alternative names
* Mutual Trust: I trust you, You trust me.
* Knowledge Mesh: Workflow knowledge is decentralized. It is distributed among all workers.
