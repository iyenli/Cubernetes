from serverless import Request, Response

def action(req: Request) -> Response:
    name = req.param('name')

    return Response("Greetings, traveler {} from beyond the fog...\n".format(name))