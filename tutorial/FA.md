# آموزش

1. ابتدا فایل مورد نیاز از بخش Releases دانلود کنید

![Img1](./img/1.png)

![Img2](./img/2.png)

2. فایل زیپ دانلود شده را از حالت فشرده شده خارج کنید

3. فایل conf را با نوتپد باز کنید

![Img3](./img/3.png)

![Img4](./img/4.png)

4. فایل کانفیگ با توجه به توضحات اصلاح کنید
    * `Hostname`: دامنه یا ساب دامنه ای که در کلودفلیر ثبت کرده اید با پروکسی روشن
    * `SNI`: دامنه شما بدون ساب دامنه
    * `MaxPing`: حداکثر پینگ مجاز
    * `Maxlatency`: حداکثر تاخیر مجاز برا اساس درخواست HTTP
    * `Jitter`: فعال یا غیرفعال سازی محاسبه جیتر.
    * `MaxJitter`: حداکثر جیتر مجاز.

5. فایل cfscanner را اجرا کنید

6. خروجی ها در فایل `result.txt` قرار خواهد گرفت

## نکات مهم

* شما میتوانید از اسکنر درحال اجرا خارج بشوید و همچنان خروجی ها را در فایل `result.exe` داشته باشید
* گزینه `Goroutines` در فایل کانفیگ برای تعداد اسکن های همزمان میباشد. با افزایش این مقدار اسکن سریع تری خواهید داشت ولی درصد خطا تاخیر افزایش میابد
