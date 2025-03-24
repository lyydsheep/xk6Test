import http from "k6/http";
import alipay_signer from "k6/x/alipaySigner";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";
import { check } from "k6";

export function isFieldNotEmpty(obj, field) {
  return obj.hasOwnProperty(field) && obj[field] !== null && obj[field] !== "";
}

export default function () {
  const gateway = "http://47.86.177.131:18881";
  const paymentRequestId = uuidv4();
  const path = `/api/v1/payments/pay/${paymentRequestId}`;
  const method = "POST";
  const privateKey = `MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCzX5NXYYzxnesmLmkoKXcFRju0k2T9vty7k8uN1yUnXpXxqqKSZwOGu6Dd74loW0GDkHTxo/ax0RQ2EbwqVcy6M5k+MOueKvPyN8EoVRTrt2HTlHjeP6TEFhc36IWr8/zV73lQ122JWjB8VDoDwpSfRNq9FwQQwnSM+tUluil7bFlDf2XwWq/RUnErk3ucpY47RA8tkLrsVqiG8ezop8UU+Cj8F838M37t5g1Ws0n/Iu016GQvCvANOMJXa0LKlkCuF6Qw/psGgAtukY/Dpa8WKKRZZRjRmeGd/bpVUjHbMoiveut+1oxXq1/ViuajcmeHB+dUjEOLFF6HM3rOUTaPAgMBAAECggEAIFtHSuXv9M3V00df9Ly2GZ93ubU07Ij3bGbWuzlqHFX1mmh7dwlaG33RIPfKw6ZihQcT8Vrwb1cV3EMKvGBJ0/Tm2c8dUaTR6ImiMFJYucSUwwPvYbf9UtnhSlaJdmFG5JiUO78ApVw9js/qvy7kfG6BPsbeFt/dAmlk9a9FOFwft3IXN1diXRd9TieHHztISXFfKrogXD46UcTNXh5FI/iWE9/Gw1KD1TvHVIcuyYeGPH3Nz1g/EF4IFXFXu9gjQk+GjHRM33kLWXTQScKQVAeG3nbWU61qWjNFm/ZcmH/StBXFlJae8L6ZZVvPlfye3VeSgPd95sa8U6TJw+LZcQKBgQDioapxDdY4/l06uZhAqbq0//Q5TSYRdBxMUGkvXIVYyEAcpegigeOZb96QloNTCJkgpMQPVVKFIW5sLTxWZeUjbh2bi5bLJ4dX8cAJiPkjjSW+sskIdoAEwdO0WcEBa9rzPxBgfwwZv63OQXK01U+VGq2AHIcMa2zMQobY1W3rgwKBgQDKnieeie/2tJnzq5uTf+EzeNyIkVl88J1PLVQNTYRykwP3zFYoZpQrfDf7NOTsxR2pUCJVLAy7ShgF972Z80+7RmmqfhsVNretJ1KiXiAtq+vsfmzXAOTO91O78NO+QOfdEwgQutagsDXS+eNWRp4PqcGGhQYmOjByjsuci0JfBQKBgBl3GPPDHkMhMdCbciQx7izQZdzacmCbr2JT1r3fo4wqVCnj6oWWGsDu9Q9CGleGK86jNPSUHcWf0AyPuKvsnyawBNupf7QsKOUU4QMxRO1dutQGutgcmJ3wOZ2WSD2kpOGYQHrXS8DI7Pq0F/OB1INoj/5JNlOK7pq1DvnmYYeJAoGAbK2H8rFp1JnqOZjCScs9r64UG+xaY3lr5xwZCUma0Rmp9y/Sxri+oNRv8n3cjGLuFfK1d5m4+nwzhn/rYrfu/DQ4WQpq3GYM/wMof46dE+IzGRZ2qpwAHkLq1tPFvzZxJ1Md8FtG48mgFRmTpqMaBKy48L5JHhf4BHozHDRV1UECgYA56lJCUWCpMXpOURFu7UCvd7RcTNOH1Wx3P791BxHFg3/f/16OtWoqyG0go2u4IPWc/wV18h0AldqPWXYDr59QtcqEjjHyAEq0T+5wyjauknhxlfYuB+qYx1UDmT3YnJatIoi82h/6J4VHSx6TsNyHCMMsuqwLiAT1iuyl0X442g==`;
  const clientId = "its4test-test-test-test-ittruly4test";
  const url = `${gateway}${path}`;
  const jsonReq = JSON.stringify({
    env: {
      terminalType: "WEB",
      clientIp: "140.205.11.6",
      deviceTokenId: "UtmF4tnJ3rWptMXWzBrko3RPiZ8uKkNYtpiXPk9eG5WOnPM1lQEAAA==",
      osType: "ANDROID",
    },
    order: {
      referenceOrderId: "2b7daafe9-95b0-44d5-ab11-56ef034119d4",
      orderDescription: "example order",
      orderAmount: {
        currency: "USD",
        value: 100,
      },
      buyer: {
        referenceBuyerId: "cidtest_notify_session01",
        buyerName: {
          fullName: "fullName",
          firstName: "firstName",
          lastName: "lastName",
        },
        buyerEmail: "test@xx.com",
        buyerRegistrationTime: "2019-11-27T12:01:01+08:00",
      },
      goods: {
        referenceGoodsId: "11451qwerqtwrq2wr",
        goodsName: "GoodsName",
        goodsQuantity: 1,
        deliveryMethodType: "DIGITAL",
        goodsUnitAmount: {
          currency: "USD",
          value: 100,
        },
        goodsCategory: "Hosting",
      },
    },
    paymentMethod: {
      paymentMethodType: "CARD",
    },
    settlementStrategy: {
      settlementCurrency: "USD",
    },
    paymentAmount: {
      currency: "USD",
      value: 100,
    },
    paymentNotifyUrl: "https://close-disadvantage.com/",
    paymentRedirectUrl: "https://spirited-passport.name/",
    productCode: "86",
    availablePaymentMethod: {
      paymentMethodTypeList: [
        {
          paymentMethodType: "CARD",
        },
        {
          paymentMethodType: "ALIPAY_CN",
        },
      ],
    },
    subscription: {
      periodType: "MONTH",
      activeTime: "2025-11-27T12:01:01+08:00",
    },
    productScene: "CHECKOUT_PAYMENT",
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



