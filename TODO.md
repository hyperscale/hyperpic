TODO
====

* Detect if image source is alpha, if yes use format supported alpha (webp if accepted or png)
* Create inteligent compression algo and choose best format by context
* Adjust quality by dpr. Ex: w=400&dpr=1 => quality=75 and w=400&dpr=2 => quality=55
* Add metrics (with prometheus or influxdb) bytes send, nb request (by source and cache), time request...
* Try use jpgoptim and other tool for optimize cache file.
* Add cache cleaner by hit or access time.
* Add endpoint for upload image with token.
* Add face detect (opencv ?)
* Fix crop region (x, y)
* Add other crop type (top-left, ...)
* Add watermark
* Add preset support by file config. Ex: my-preset.json
* For speed use small image for create other small crop and not the original image.
* Add S3 source provider
* Add Azure Blob source provider
* Add Ceph source provider
* Add Hot cache for best images (memory cache provider ?)
* Add cluster mode (???) or use DB distributed (own ??)
* Add UI (???)

Articles
--------

* https://blog.imgix.com/2016/03/11/auto-compress.html
* https://blog.imgix.com/2016/03/30/dpr-quality.html
* https://blog.imgix.com/2016/03/09/save-data-client-hint.html
