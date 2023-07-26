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
