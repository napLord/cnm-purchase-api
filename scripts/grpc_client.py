import asyncio

from grpclib.client import Channel

from ozonmp.cnm_purchase_api.v1.cnm_purchase_api_grpc import CnmPurchaseApiServiceStub
from ozonmp.cnm_purchase_api.v1.cnm_purchase_api_pb2 import DescribePurchaseV1Request

async def main():
    async with Channel('127.0.0.1', 8082) as channel:
        client = CnmPurchaseApiServiceStub(channel)

        req = DescribePurchaseV1Request(purchase_id=1)
        reply = await client.DescribePurchaseV1(req)
        print(reply.message)


if __name__ == '__main__':
    asyncio.run(main())
