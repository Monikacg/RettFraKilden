%% Plot of states matrix
%load states.mat

hold off

time_plot_102 = states_measured(1,:);
travel_plot_102 = states_measured(2,:);
travel_rate_plot_102 = states_measured(3,:);
pitch_plot_102 = states_measured(4,:);
pitch_rate_plot_102 = states_measured(5,:);
elevation_plot_102 = states_measured(6,:);
elevation_rate_plot_102 = states_measured(7,:);

subplot(711)
stairs(u_star(:,1),u_star(:,2)),grid
ylabel('u')
subplot(712)
plot(time_plot_102,travel_plot_102,'m',time_plot_102,travel_plot_102,'m','LineWidth',2),grid
ylabel('Travel')
subplot(713)
plot(time_plot_102,travel_rate_plot_102,'m',time_plot_102,travel_rate_plot_102','m','LineWidth',2),grid
ylabel('Travel rate')
subplot(714)
plot(time_plot_102,pitch_plot_102,'m',time_plot_102,pitch_plot_102','m','LineWidth',2),grid
ylabel('Pitch')
subplot(715)
plot(time_plot_102,pitch_rate_plot_102,'m',time_plot_102,pitch_rate_plot_102','m','LineWidth',2),grid
ylabel('Pitch rate')
subplot(716)
plot(time_plot_102,elevation_plot_102,'m',time_plot_102,elevation_plot_102','m','LineWidth',2),grid
ylabel('Elevation')
subplot(717)
plot(time_plot_102,elevation_rate_plot_102,'m',time_plot_102,elevation_rate_plot_102','m','LineWidth',2),grid
xlabel('tid (s)'),ylabel('Elevation rate')