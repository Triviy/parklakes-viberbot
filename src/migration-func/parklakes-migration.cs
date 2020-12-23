using System;
using System.Net.Http;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;

namespace Parklakes.Migration
{
    public static class parklakes_migration
    {
        private static HttpClient _client = new HttpClient();

        [FunctionName("parklakes_migration")]
        public static async Task Run([TimerTrigger("0 0 1 * * *")]TimerInfo myTimer, ILogger log)
        {
            log.LogInformation("starting car owners migration job");
            var apiKey = Environment.GetEnvironmentVariable("API_KEY");
            var apiURL = Environment.GetEnvironmentVariable("API_URL");
            try {
                var request = new HttpRequestMessage {
                    RequestUri = new Uri($"{apiURL}/api/v1/car-owners/migrate"),
                    Method = HttpMethod.Post,
                    Headers = {
                        { "X-API-KEY", apiKey }
                    }
                };
                var result = await _client.SendAsync(request);
                log.LogInformation($"ending car owners migration job with status code: {result.StatusCode}");
            } catch(Exception e) {
                log.LogError(e, "calling migration API failed");
            }
        }
    }
}
