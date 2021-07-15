![aaBanner](https://user-images.githubusercontent.com/15215359/121749073-0de83300-cad8-11eb-8306-d8f69334731b.jpg)


A desktop notification app for Apex Legends map rotations. Using https://apexlegendsapi.com/ and [Beeep](https://github.com/gen2brain/beeep), a cross-platform Go notification library.

If you wish you use this, an API key can be generated (relatively) anonymously at apexlegendsapi.com.

## Tested Platforms

##### Linux (ArchLinux) running Budgie DE (GTK-based):

![image](https://user-images.githubusercontent.com/15215359/121745794-dd51ca80-cad2-11eb-8144-dfc0d3a39c4e.png)


##### Windows 10 (in a qemu/kvm VM):

![image](https://user-images.githubusercontent.com/15215359/121760302-02a4ff80-caf8-11eb-9f43-43689ecb37a4.png)


As I test on new environments, I will update the README. It's obviously a bit of a kludge to have a cmd window up in Windows (and to run from a shell in \*Nix as well), so I may look into running this as a tray applet (on top of all the other things I have planned). On top of making for cleaner execution, this would also have the benefit of allowing UI-driven configuration, and more stylized notifications (Map-as-background, for instance). That's kinda a re-write though, and this was a simple project to mess around in Go a bit more.
