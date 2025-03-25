import http from "k6/http";
import alipay_signer from "k6/x/alipaySigner";
import { check } from "k6";

export function isFieldNotEmpty(obj, field) {
    return obj.hasOwnProperty(field) && obj[field] !== null && obj[field] !== "";
}

export default function () {
    const gateway = "http://47.86.177.131:17070";
    const paymentRequestId = Math.random();
    const path = `/ams/api/v1/payments/pay`;
    const method = "POST";
    const privateKey = `MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCKks+lpbs5JZ5p
QcyCVvF2PGJE49UNS3maxq7qTCbFAF7EpMj3lHjlRa70GfCwtyhAYs81PpYs5nRD
1iD+U7RfUVuNqPawGkm3BKYrLNYBbxcsTDZ/OoAtXINsPXl855s2jTY5EgQO+fXJ
YfVjUsKeh6w4IZdHzpNMgyyOYiz0w3bNx3aNsbZ6fewbVg6z5GmTg0h9ltLQ5qMK
DbPPmZJ12m6ebDvK9BIf3I9s5aP36K1Rv2xXB6WRJtT5p7BkLXt9qE7iCGr8gV/x
3ksU3zKb8MHsxyfDWD544rb3tPPPsy3p1EE+TA4ro8RZRayTGgjnTULF2Clbex6a
PrA3+zthAgMBAAECggEACJgJdkjZXAMjPq/f+Qb47pylrSDcEHkDAtnmwY0VEnHV
bsIxdhLbjiYZYlyBetQneoBZc80WQmHNPOL4vnN+fWYHex+X1pARMxJ/8Hy7wPCc
y7IZgLmmEJCEoTlzlI8TwMXwYojlfr0UKysYegWiudlikUPNpM+afv8Oy/D2v/v2
WSZEqGikwaffE2MvyA8ss3gVWCvmyNv89/Q/goszy+PGoV5EV/yZ/RwbbznNs6Pb
SnXQ9J0j1tRY0L8tNizDnoElRRJ2h8a+6iogbdKUVAoPiDiwI5nWWdN8wCtrWmfR
pq7vbRlhXGrrX6FWS3qy6patvwmXNXXQL1IQjEk9aQKBgQC/kcf5pEf4t3bBHnlP
M4dHeFqAhfQwFPp6WBd7Wl/PM6PInqcuBEZ5FSCYB9Bh+sqqoWZ+pa0H6DT1OGoW
XBobSrogy3YXLSSPbkgliMMYUpnn5dQnAyW8LTeJYZ6q1gtJ+qjsI09W+Rdalf91
/1R+3JJS1HSGY3g2XffpGeu1nwKBgQC5Lg1ICVpma1ZVPCjC7Mpg5nqWHecUNkSg
kth/71K1e+998Y3SGY72h/6fIW8K35Edhq6GfqgpKcB9gOULzL6VocM844T5ZcSl
uh37hrrxlXPAXIu703KBdc3ItFW0Q417rGrqh/PRdnRw9X4GAF9TcLEw7+TktdZC
iUeq+nlu/wKBgH2Z/mxeWtXmrBT4fv7/wO2KKoRjz3OK/aMjiNnWqkS3DeampuQT
54TR5lnnnafv/9saEZJt2+H4TGiPQXdBkdhdCWYhIF8XuQXVf7YkUg3rcn9J/+xI
MwCLAQOxHo/R4PrzPrf8erOCg95fxGvAKc03nzRxmajXJOU4fSe3WATvAoGBAK98
9/MrgfMLh45Q66QGOKfp44Q9pE5gO1scrnVXPL9mSwjEkIzp0bTKHj95tLzOL7yW
dPBaOUyBF70YGHe9OWOeH+KlDtA2ZExV+7Hw9VqaMk66pWWDNcF//VtVot6pIfxw
4gWOfz4ijqi5zQss8Smm4xSoUvd3Zyw44qUip0/LAoGBAJhZf3V0XSpE8KRdgzfu
7odc5BdS9UGswZZke17fz0INLs1Saa8fL46gQ6Quqb5jpD3dS4FdFh2daURDnykv
nh5GhYcp3qp29ppExQscuy0e276k8fiBSo2wwiMcFI/tNBKdQBoN/F+NZ8cwMIpB
cgSnHSATLKivhz9T/E6JFJ25`
    const alipayAlipayPublicKey = `MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0xPY0DkRrOOPGlEJyXop
nqPPWGLYB/hUeo59PUGLzBYXJ4ktM0KZY5szT+D9zz2IrfIJ+xarl+7J1Pb+kfwr
pXcmJ1h/IAhsCFEdfiaBHrCMtesazZ7l09D81IoNYUQ6nrdfCjO93LpAqt8eLFKH
ujqHAafXW3ww9nasua5dmF5kXBMMfESRsPVVGjfJTf3wcRHJ8GnIwS3X1b8sCcT/
GZ4icAMRxcCnihtWty/NVf7H/fvSVU8LzOcZCrboQvxq+an+lKRmx9KOLgOd040e
pCChRMCn00m/ulpaTk8Hh6Wd6IJx08kfsfp6sRCfxtlf9+X9B73hdwE9gSB4iIoV
3QIDAQAB`

    const clientId = "391383a3-ff20-11ef-98d4-00163e233080";
    const url = `${gateway}${path}`;
    const jsonReq = JSON.stringify({
        "productCode": "CASHIER_PAYMENT",
        "paymentRequestId": paymentRequestId,
        "order": {
            "referenceOrderId": "228",
            "orderDescription": "ClawCloud Staging - 账单 #228",
            "orderAmount": {
                "value": "388",
                "currency": "USD"
            },
            "goods": [
                {
                    "referenceGoodsId": "pid-2",
                    "goodsName": "VPS - 1C / 512M / 10G / 500G",
                    "goodsCategory": "Hosting",
                    "goodsUnitAmount": {
                        "value": "388",
                        "currency": "USD"
                    },
                    "goodsQuantity": "1",
                    "deliveryMethodType": "DIGITAL"
                }
            ],
            "buyer": {
                "referenceBuyerId": "2",
                "buyerName": {
                    "firstName": "JIE",
                    "lastName": "YANG",
                    "fullName": "JIE YANG"
                },
                "buyerEmail": "yangjie198912112@gmail.com",
                "buyerRegistrationTime": "2024-12-11T08:01:23+00:00"
            }
        },
        "paymentAmount": {
            "value": "388",
            "currency": "USD"
        },
        "paymentMethod": {
            "paymentMethodType": "CARD",
            "paymentMethodMetaData": {
                "billingAddress": {
                    "region": "CN"
                },
                "is3DSAuthentication": false,
                "tokenize": true
            }
        },
        "paymentRedirectUrl": "http://www.yourRedirectUrl.com",
        "paymentNotifyUrl": "http://47.86.177.131:17070/test",
        "paymentFactor": {
            "captureMode": "AUTOMATIC",
            "isAuthorization": true
        },
        "settlementStrategy": {
            "settlementCurrency": "USD"
        },
        "env": {
            "terminalType": "WEB",
            "userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
            "deviceTokenId": "UtmF4tnJ3rWptMXWzBrko3RPiZ8uKkNYtpiXPk9eG5WOnPM1lQEAAA==",
            "clientIp": "140.205.11.6"
        }
    });

    var headers = alipay_signer.genSignatureHeader(
        clientId,
        path,
        method,
        privateKey,
        jsonReq
    );
    const param = {
        headers: headers,
    };
    const res = http.post(url, jsonReq, param);
    console.log(res.body);
    console.log(res.status);

    check(res, {
        "status is 200": (r) => r.status === 200,
        paymentSessionCreated: (r) => isFieldNotEmpty(r.json(), "normalUrl"),
        statusSuccessful:(r)=> r.json().result.resultStatus==="S"
    });
}



