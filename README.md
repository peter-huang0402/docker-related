<b> Research for live migration in docker</b>

The Project is related to live migration in docker between two different machine.
We update the original patches and enhance and expand its functionalities with 3 approaches.
However, currently the sucessful live migration in docker just is limited under some specific cases in containers. 
Under the specific cases in containers the live migration works well with 3 approaches.

<b>(1) Stop and Resume approach</b>
    Dump all memory and docker-related data into dumping file. Based on dumping file, the container can be restored on another machine. 
     
<b>(2) Pre-Copy approach</b>
    Just dump the updated memory pages only. Each time the dumping action compares with the previous dumping page, and just dump updated memory pages. Finally, we combine all dumping pages to restore the containers.

<b>(3) Post-Copy approach</b>
    Just dump the initial pages to another machine. When page faults happen in destination machine. It will make requests of missing page from the destination machine to the source machine. And, the source machine will send the requested pages back to destination machine.

