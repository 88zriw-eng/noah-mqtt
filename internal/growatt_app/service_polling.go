func (g *GrowattAppService) pollStatus(device models.NoahDevicePayload) {
	if data, err := g.client.GetNoahStatus(device.Serial); err != nil {
		slog.Error(
			"could not get device data",
			slog.String("error", err.Error()),
			slog.String("device", device.Serial),
		)
	} else {

		payload := devicePayload(data)

		last := g.lastGenerationTotal[device.Serial]

		// Ungültige oder zurückspringende Gesamtenergie ignorieren
		if last > 0 {
			if payload.GenerationTotalEnergy <= 0 {
				slog.Warn(
					"Ignoring invalid GenerationTotalEnergy",
					slog.String("device", device.Serial),
					slog.Float64("received", payload.GenerationTotalEnergy),
					slog.Float64("previous", last),
				)

				payload.GenerationTotalEnergy = last

			} else if payload.GenerationTotalEnergy < last {

				slog.Warn(
					"Ignoring decreasing GenerationTotalEnergy",
					slog.String("device", device.Serial),
					slog.Float64("received", payload.GenerationTotalEnergy),
					slog.Float64("previous", last),
				)

				payload.GenerationTotalEnergy = last
			}
		}

		g.lastGenerationTotal[device.Serial] = payload.GenerationTotalEnergy

		for _, e := range g.endpoints {
			e.PublishDeviceStatus(device, payload)
		}
	}
}
