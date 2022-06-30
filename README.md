# ddcpopup

This application is a part of the [nwg-shell](https://github.com/nwg-piotr/nwg-shell) project.

`ddcpopup` allows to control basic settings of external (not laptop built-in) monitors. It's a simple frontend to [ddcutil](http://www.ddcutil.com). 
It provides 2 features:

1. creates a GUI for brightness, contrast and monitor preset settings;
2. returns textual output for use with panels/bars as a brightness widget:
- 2 lines long output for nwg-panel / Tint2: icon path / brightness percentage, or
- 1 line long brightness percentage with a user-defined label.

![2022-06-30-012005_screenshot](https://user-images.githubusercontent.com/20579136/176566079-7ae68b33-ec67-4ac4-8eeb-205d46ac7af2.png)

**This software is aimed at power users, rather then beginners, as it may turn out quite difficult to set up.**

The `ddcpopup` command is being developed with the [sway](https://github.com/swaywm/sway) Wayland Compositor in mind, but should work in other environment, including X11. The latters has not yet been tested, however.

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

## Dependencies

- `ddcutil`
- `go` >= 1.16 (make)
- `gtk3`
- `gtk-layer-shell`

## Installation

1. Install [ddcutil](http://www.ddcutil.com), give it all the necessary [permissions](https://www.ddcutil.com/i2c_permissions). You may also need 
to start `i2c-dev` in `/etc/modules-load.d/i2c-dev.conf` (or another location named in `man modules-load.d`). Check the `ddcutil environment` command output if something goes wrong.

2. Clone this repository, cd into it:

```text
git clone https://github.com/nwg-piotr/ddcpopup.git
cd ddcpopup
```

3. Build and install

```text
make build
sudo make install
```

4. Find your `/dev/i3c-` bus numbers to use, e.g. `4` and `5` in the output below:

```text
$ ddcutil detect
Display 1
   I2C bus:  /dev/i2c-4
   EDID synopsis:
      Mfg id:               AOC
      Model:                22V2WG5
      Product code:         8706
      Serial number:        
      Binary serial number: 116 (0x00000074)
      Manufacture year:     2020,  Week: 2
   VCP version:         2.2

Display 2
   I2C bus:  /dev/i2c-5
   EDID synopsis:
      Mfg id:               AOC
      Model:                2475WR
      Product code:         9333
      Serial number:        
      Binary serial number: 16843009 (0x01010101)
      Manufacture year:     2016,  Week: 47
   VCP version:         2.1

Invalid display
   I2C bus:  /dev/i2c-7
   EDID synopsis:
      Mfg id:               AUO
      Model:                
      Product code:         53485
      Serial number:        
      Binary serial number: 0 (0x00000000)
      Manufacture year:     2019,  Week: 11
   DDC communication failed
   This is an eDP laptop display. Laptop displays do not support DDC/CI.
```

## Usage

### nwg-panel executor

![image](https://user-images.githubusercontent.com/20579136/176570813-505eb8d7-de10-4d57-9e7b-b4056b25853d.png)

Sample nwg-panel executor for the monitor on bus #4:

- Script: `ddcpopup -b 4 -e`
- On left click: `ddcpopup -b 4 -hpos right -vpos bottom -hm 6  -vm 6` (popup at bottom-right, both margins 6 px)
- Interval: `5` seconds. **Lower values may interfere with the popup widow, as it sends 4 DDC/CI requests, one after another. Each takes at least 300 ms. It's impossible to execute two requests at a time**

![image](https://user-images.githubusercontent.com/20579136/176569246-67992f4d-af91-470e-b1fc-a83e2cc3902e.png)

### nwg-panel button

A safer and more resources-friendly usage may be just a button:

![image](https://user-images.githubusercontent.com/20579136/176570949-83034f31-2b24-4a9b-b4ad-b331b233b5d7.png)

![image](https://user-images.githubusercontent.com/20579136/176571080-539ec985-d671-4992-bb7c-7008d9996405.png)

