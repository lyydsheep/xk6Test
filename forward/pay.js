import http from "k6/http";
import alipay_signer from "k6/x/alipaySigner";
import { check } from "k6";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";

export function isFieldNotEmpty(obj, field) {
    return obj.hasOwnProperty(field) && obj[field] !== null && obj[field] !== "";
}

export let options = {
    stages: [
        { duration: '30s', target: 50 },
        { duration: '1m', target: 100 },
        { duration: '30s', target: 50 },
        { duration: '30s', target: 0}
    ]
}

export default function () {
    const gateway = "http://47.86.177.131:17070";
    const paymentRequestId = uuidv4();
    const path = `/ams/api/v1/payments/pay`;
    const method = "POST";
    const alipayMerchantPrivateKey = `MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDEiH/PTT+RWYi7
+xV9QeSrKW6E/xCgedqdgBI6qbUvEhAqkZzOmRh4KhNInthi9r36l8gKiA2VJ4lr
K06yUfThgY0nfZ6AKT64V4rbHsx0HFrI/NKRXo61c3406Ul/Oqnvsmd5GrU6xb7+
uBZA8QbFE7ToslvwiAfvo9RpbRoStmT0vx2PvudNbt/gBTYtv2waTYaPj95bZyvR
9Yo96dGcC1UotNfq4vt/3JXjcgH3x26UCDwoeYNcP1snDB1HW5T6WVc7xW8W8Yv0
z0yARSEd48utg8xc59xdpyECkYPlROvw/7lVrZRCOrl/agxXWtKrks8JeTB8v8d0
a9KQ/gb5AgMBAAECggEAQX23EY6NY1MxtGqsM4xUthDdamZQ1RkmF5wI9xF8dg4P
4w+Y8GOl+s0Slf2Q4BGXJz2TBKdn109QhKzu0Y9sCgWQ0xsSSWT1KJrLg89UlFCq
lBVj2dtntOGiqeEfg805uc16m6mhwM8KoXkYuVVYSy4Q+AYPiyzZcCro6qTXrmoR
g8Cg+/in1ItHowyfSOLzRsqvbp9vGdU3lLrWi5MmXLmTTalotJHXc0Rd9Slrh2bY
3YubhrZfDZlApmkSS9rTLoED+tI2AZdI49W6pnjHB6AfBgCXVsCK7yjitj/iKDqC
reNdhK5eY3Cwn1RErEJGOEM0VyGaeb1ut+nC65ldgQKBgQDw+e9RjCK8ZVec8VzB
uU7OFxUBcAmNVOJ1jx7b3SStq5wmuPLW3dtbPN8V0kBFf/T0pfaBgKj79S5miMmw
ERM9C2FenghN46D9uXVoxFqHNKe1kn+2dpuRkAd2IeRsrWn2XDeU9DXYM4fZ0itq
2BPJfYpOdoOnta37B56goiiFkQKBgQDQyTyv4rjOoKSg4aS2LMqYM5A49fKh5Lg/
2DWCl0gZ2Q+iQqSEwt57Jejj3QO6+z8aMMlDD/qROeKvcCLWv1Ym2VbGM3yr7zSg
t108hudiKuOBlZ9XZiEiioc6cu0owNIfM6zNIdpaKwG9V6Oa/DZ3g5yUwXMnR82g
/Yp/e4IW6QKBgQCu0SLEzh0E/6AnwxG/mGeLK0OZ32WOml4PWtzQNAY/15dYoCPL
rPdNoUNV2Um3IbTbJutF18i/wIcA64slp72FM5RXx93OY6yPZNPARXJHU/O2zajI
/hKt7wb6tGu6S7Prfcr0zJWjWv7bDpVg1ZDFQ8XqVh/8sticnFJ/xiQPgQKBgHTp
Wrw6vrWlqsoT0EHazw9vQEFFJ7qT8sB9d2lLASrIK0L3Alz9KcvXrJN7/UzEx88I
poqQ9gRAX7lRl5Ccz8ctSLPvvM4iQlwEkYcFG6gS0BaODA3KuJ845wRJupdpcb/b
FdZAMJ7xGiZGXuy4cl92KUX7FVpXkMOnddhw9qWRAoGAAvV4Gjd/E2FPWgHAwtvl
vmP76FeVbr9rs6/pl+GwEQ1f6ghI1ZEX8L0jUD9YAgecHj0Dw21CgfvkOBoitkn4
xw2AynE1bD8oQZd8C3VYf/8h/BOBgr5nDWsa6dXl6+ceb1S7+0P1VQZ+pwLkElYI
AfYtwrx5Qncml9u9IHDxxWY=`


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
        "paymentNotifyUrl": "http://www.yourRedirectUrl.com",
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
        alipayMerchantPrivateKey,
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
        checkStatus:(r) => r.json().result.ResultStatus === "U",
        checkMessage: (r) => r.json().result.ResultMessage === "payment in process"
    });
}



