go-aurora
=========

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

There is also a "helper" function called `Check()` which takes the results of `Get()`, checks the Speed, Density and Bz (the three measurements aurora fans are most interested in) and gives you a value between -1 and 2 which indicates possible aurora activity:

    -1 = ACE spacecraft returned -999.9 for Bz, Density and Speed. This indicates no usable data
    0 = Bz is greater than the warning threshold, while Speed and Density and below the warning threshold
    1 = Bz is less than the warning threshold, but greater than the alert threshold. Speed and density are greater than their respective warning thresholds, but less than their alert thresholds
    2 = Bz is less than the alert threshold, and Speed and Density are greater than the alert thresholds. Grab your camera and run!

Helping out
===========

I've only just become interested in auroras (and only captured my first one in August 2015), so I might be

To-Do:
======

- [ ] Look for Aurora Australis data, as this data may be for Aurora Borealis only
- [ ] Modify the code so it grabs all of the data, instead of just the last (latest) line
- [ ] More detail in the Check() func, as 2/3 data points could be bad, but still reports as 0
- [x] Add Kp data from ftp://ftp.swpc.noaa.gov:21/pub/lists//wingkp/wingkp_list.txt (or similar)
