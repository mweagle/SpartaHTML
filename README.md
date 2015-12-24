# SpartaHTML
[Sparta](https://github.com/mweagle/Sparta) application that demonstrates provisioning an S3 backed site with a CORS-enabled API Gateway

<div align="center"><img src="https://raw.githubusercontent.com/mweagle/SpartaHTML/master/Sparta.gif" />
</div>


## Instructions

1. `git clone git clone https://github.com/mweagle/SpartaHTML`
1. `cd SpartaHTML`
1. `make get`
1. `S3_BUCKET=<MY_S3_BUCKET_NAME> make provision`
1. In the _Stack output_ section of the log, look for the **S3SiteURL** key and open the provided URL in your browser (eg: _http://spartahtml-site09b75dfd6a3e4d7e2167f6eca73957e-zp9okcokn7o.s3-website-us-west-2.amazonaws.com_).

## Result

<div align="center"><img src="https://raw.githubusercontent.com/mweagle/SpartaHTML/master/site/screenshot.jpg" />
</div>
