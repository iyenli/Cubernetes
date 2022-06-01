from serverless import Request, Response

def action(req: Request) -> Response:
    name = req.param('name')

    return Response("Hello, {}!\n".format(name))
