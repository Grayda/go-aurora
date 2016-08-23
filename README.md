go-aurora
=========

**2016 Update: If you want an easier way to get this data, check out http://developers.auroras.live . All the heavy lifting is done for you, and you just need to parse a simple JSON. Easy!**

A simple library to help astrophotographers capture the Auroras. I wrote it so I could use it with my Ninja Sphere and get alerted whenever there was potential aurora activity

Usage
=====

Simply import the library, then call:

    results := aurora.Get()

`results` is a map of strings that contains your data. Supported data types include:

    YR (UTC Year)
    MO (UTC Month)
    DA (UTC Day)
    S (Status. 0 = Good data, 1-8 = Bad data, 9 = No Data)
    Time (HHMM)
    Seconds (Seconds of the day)
    JulianDay (Modified Julian Day)
    IonTemperature
    Speed (Bulk Speed)
    Density (Proton Density)
    Bx
    By
    Bz
    Bt
    Lat
    Long
    Kp
    Kp1hour
    Kp4hours

There is also a "helper" function called `Check()` which takes the results of `Get()`, checks the Speed, Density, Bz and Kp (the four measurements I think aurora fans are most interested in) and gives you a score out of 90, on how likely you are to see an aurora. The scores are weighted, with Kp having less significance in the final score. Weights are currently;

    KpGreenWeight  = -5
    KpYellowWeight = 2
    KpOrangeWeight = 5
    KpRedWeight    = 10

    GreenWeight  = -10
    YellowWeight = 10
    OrangeWeight = 15
    RedWeight    = 25

Essentially, if Bz is in the "green" (between 20 and 0), 10 points are taken off. If Density is in the red, add 25, and so forth. I'm terrible at mathematics, statistics and so on, so if you can come up with a better model, PLEASE put in a pull request or message me!

Helping out
===========

I've only just become interested in auroras (and only captured my first one in August 2015), so the data contained here might not be correct, or I might be grabbing the wrong data, or not enough data. I strongly encourage anyone with more experience in aurora hunting to suggest changes, either code-wise, or data-wise. You don't have to be a programmer to help!

To-Do:
======

- [ ] Look for Aurora Australis data, as this data may be for Aurora Borealis only
- [ ] Modify the code so it grabs all of the data, instead of just the last (latest) line
- [ ] More detail in the Check() func, as 2/3 data points could be bad, but still reports as 0
- [x] Add Kp data from ftp://ftp.swpc.noaa.gov:21/pub/lists//wingkp/wingkp_list.txt (or similar)
