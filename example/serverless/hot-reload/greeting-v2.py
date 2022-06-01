from serverless import Request, Response

def action(req: Request) -> Response:
    name = req.param('name')

    return Response("{}, thou shouldst take the crown...\n".format(name))