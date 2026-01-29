# Installation steps

## Step 1:
Download the STM32 Cube IDE to run code on the microcontroller

Download the code from this repostiory https://github.com/vanDeventer/as2/tree/adc and run it in the IDE for the microcontroller

## Step 2:
Connect the microcontroller to the raspberry pi and check which serial port is used by the microcontroller. Usually for raspberry pi 5 it is the "/dev/ttyACM0" directory

If you are not using raspberry pi (Linux) and is using windows then check the serial port by checking the device manager. Go to COM ports to see which one is the microcontrollers UART

## Step 3:
Clone this repository and download GOLang

### How to download GOLang:
On raspberry pi 5 (linux) follow the steps on this site https://go.dev/doc/install

On windows follow the steps (for windows) on this site https://go.dev/doc/install

### Change the systemconfig.json file:
Check on the traits and see that the port is correct. I.e. the serial port from step 2

### How to run the code:
Go into the potentiometer folder from the terminal and use the command "go run ." (without the "") to run the code

## Step 4 (Optional only to see if it works):
Copy this URL to the web: 130.240.153.92:20151/rudder/doc

When running the code you should see a json file with the value of the potentiometer
