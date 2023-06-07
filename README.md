# DragonTable
 
This is a golang based application that was made specifically for a 42" Elo Touchscreen Display that is embedded in the surface of a table made for table top gaming. Since this was made for a very specific device you will see interops to a C++ library for the Elo drivers. The app should be usable for any touch screen device if you remove the references to the Elo drivers. If you are using this on a different touch screen device the windows touch drivers should natively allow this to work just fine. While developing this app my goal was to see how far I could stretch Golang for a user interface experience. 

This app load and displays high resolution images (maps) on a screen from a configured folder on the device. It includes features like an overlay of a 1x1 grid, and the disabling of touch for the use of mini-figures on the screen.

Note: I have only tested this on Windows, but it should theoretically work on other operating systems with some tweaking.

