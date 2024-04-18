

https://github.com/Paintersrp/tmux-launcher/assets/82731232/27417edb-4cff-4b40-9e7e-a05658272677

The bottom portion of the window allows selecting from a list of previous tmux sessions. The top portion of the window shows a preview of the **active** pane in the highlighted session. 

# Why?

Good question. Once in Tmux there are plenty of sessionn navigation options but when coming in fresh I was always stuck awkwardly typing "tmux ls" because I'm bad at remembering what I name my sessions. Then I would have to go through the arduous journey of typing "tmux a -t {session_name}". Look at how many letters that is...

Joking aside, it's been slightly helpful in my workflow and I liked making it. That's why. 

For video context I've bound a shortcut on my OS to launch a shell which runs the compiled tmux-launcher. If you clone it, you can do the same or simply build and run the compiled file. 

## Are you going to add to it?

I will overtime improve the overall interface. I could emulate some tmux plugins that exist that allow managing (deleting, renaming, etc) sessions but realistically they already do a great job of that and the only tangle benefit would be that they could be managed without launching into tmux first. This script's primary purpose is just launching into the right session from a fresh boot to get going. 

My next plans are to include simple apis for launching into new sessions of a given layout. Like a fresh "ide" session for example which might open up a main window for a text editor and smaller shells for test running, etc. 
