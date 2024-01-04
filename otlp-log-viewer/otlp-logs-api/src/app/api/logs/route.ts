export const runtime = 'edge'

export async function GET() {
    return Response.json({
        resourceLogs: {}
    }, {
        headers: {
            'Access-Control-Allow-Origin': '*',
        }
    })
}

export async function OPTIONS() {
    return Response.json({}, {
        headers: {
            'Access-Control-Allow-Origin': '*',
        }
    })
}