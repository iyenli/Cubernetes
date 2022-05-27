from serverless import Request, Response

def action(req: Request) -> Response:
    a = req.param('a')
    b = req.param('b')

    return Response(str(a + b))