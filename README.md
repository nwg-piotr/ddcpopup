# ddcpopup

`ddcpopup` command provides 2 features:

1. creates a simple GUI, as a frontend to `ddcutil`-controlled brightness, contrast and preset settings;
2. returns output for use with panels/bars as a brightness widget:
- 2 lines long output for nwg-panel / Tint2: icon path / brightness percentage, or
- 1 line long brightness percentage with a user-defined label.

![2022-06-30-012005_screenshot](https://user-images.githubusercontent.com/20579136/176566079-7ae68b33-ec67-4ac4-8eeb-205d46ac7af2.png)

The `ddcpopup` command is being developed with the [sway](https://github.com/swaywm/sway) Wayland Compositor in mind, but should work in other environment, 
including X11. The latters has not yet been tested, however.

```text
$  ddcpopup -h
Usage of ddcpopup:
  -b int
    	Bus number for /dev/i2c-<number> (Required; check 'ddcutil 'detect') (default -1)
  -d	turn on Debug messages
  -e	print brightness Executor data
  -hm int
    	window Horizontal Margin (Wayland)
  -hpos string
    	window Horizontal POSition: 'left', 'right' or none for center (Wayland)
  -k	use darK icons
  -l string
    	print this Label instead of image path
  -v	display Version information
  -vm int
    	window Vertical Margin (Wayland)
  -vpos string
    	window Vertical POSition: 'top', 'bottom' or none for center (Wayland)
```
