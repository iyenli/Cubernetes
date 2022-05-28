from serverless import Request, Response

def action(req: Request) -> Response:
    a = int(req.param('a'))
    b = int(req.param('b'))
    print("a + b = ", a + b)

    return Response(str(a + b))
